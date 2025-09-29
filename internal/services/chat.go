package services

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"websocket-chat/internal/config"
	"websocket-chat/internal/types"
	"websocket-chat/internal/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[string]*websocket.Conn)
var clientsMu sync.Mutex
var messageChannel = make(chan types.Message)

// HandleConnections upgrades HTTP to WebSocket and authenticates via JWT from Authorization header
func HandleConnections(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Validate the token
		claims, err := utils.ValidateToken(tokenString, &cfg.JWT)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			writeJSONError(w, http.StatusUnauthorized, "email not found in token")
			return
		}
		clientsMu.Lock()
		clients[email] = conn
		clientsMu.Unlock()

		for {
			var msg types.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println(err)
				clientsMu.Lock()
				delete(clients, email)
				clientsMu.Unlock()
				return
			}

			// force sender from JWT claim to prevent spoofing
			msg.Sender = email
			messageChannel <- msg
		}
	}
}

// HandleMessages reads messages from the channel and delivers them to recipients
func HandleMessages() {
	for {
		msg := <-messageChannel

		clientsMu.Lock()
		recipientConn, ok := clients[msg.Recipient]
		clientsMu.Unlock()
		if !ok {
			fmt.Printf("Recipient %s not connected\n", msg.Recipient)
			continue
		}

		err := recipientConn.WriteJSON(map[string]string{
			"sender":  msg.Sender,
			"message": msg.Message,
		})
		if err != nil {
			fmt.Println("Error sending message:", err)
			recipientConn.Close()
			clientsMu.Lock()
			delete(clients, msg.Recipient)
			clientsMu.Unlock()
		}
	}
}

// CloseAllClients closes and removes all connected websocket clients.
func CloseAllClients() {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for email, conn := range clients {
		conn.Close()
		delete(clients, email)
	}
}
