package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websocket-chat/internal/config"
	"websocket-chat/internal/middleware"
	"websocket-chat/internal/services"
	"websocket-chat/internal/types"
	"websocket-chat/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.MustLoad()

	db, err := utils.ConnectDB(&cfg.DBConfig)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		return
	}

	if err := db.AutoMigrate(&types.User{}); err != nil {
		slog.Error("automigrate failed", slog.String("error", err.Error()))
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", services.HandleConnections(cfg)).Methods(http.MethodGet)
	r.HandleFunc("/signup", services.UserSignUp(db, cfg)).Methods(http.MethodPost)
	r.HandleFunc("/login", services.UserLogin(db, cfg)).Methods(http.MethodPost)
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"status": "ok"},
		})
	}).Methods(http.MethodGet)

	go services.HandleMessages()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: middleware.RateLimit(r),
	}

	// Start server in background
	go func() {
		fmt.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown on interrupt/terminate signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down server...")
	// Close active websocket clients first
	services.CloseAllClients()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}
	slog.Info("server exited cleanly")
}
