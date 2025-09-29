package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"websocket-chat/internal/config"
	"websocket-chat/internal/models"
	"websocket-chat/internal/types"
	auth "websocket-chat/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSignUp handles signup using email or phone
func UserSignUp(db *gorm.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Require either email or phone
		if req.Email == "" && req.Phone == "" {
			http.Error(w, "Email or phone is required", http.StatusBadRequest)
			return
		}

		if req.Password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		// Normalize (optional)
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		req.Phone = strings.TrimSpace(req.Phone)

		// Check uniqueness
		var existing models.User
		if req.Email != "" && db.Where("email = ?", req.Email).First(&existing).Error == nil {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		if req.Phone != "" && db.Where("phone = ?", req.Phone).First(&existing).Error == nil {
			http.Error(w, "User with this phone already exists", http.StatusConflict)
			return
		}

		// Hash password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:     req.Email,
			Phone:     req.Phone,
			Password:  string(hashed),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		// Generate JWT using passed config (avoid reloading config per request)
		token, _ := auth.GenerateToken(
			strconv.FormatUint(uint64(user.ID), 10),
			user.Email,
			time.Hour*24,
			&cfg.JWT,
		)

		res := types.SignupResponse{
			ID:    user.ID,
			Email: user.Email,
			Phone: user.Phone,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

// UserLogin authenticates a user by email or phone and returns a JWT
func UserLogin(db *gorm.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Password == "" || (strings.TrimSpace(req.Email) == "" && strings.TrimSpace(req.Phone) == "") {
			http.Error(w, "Email or phone and password are required", http.StatusBadRequest)
			return
		}

		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		req.Phone = strings.TrimSpace(req.Phone)

		var user models.User
		var err error
		if req.Email != "" {
			err = db.Where("email = ?", req.Email).First(&user).Error
		} else {
			err = db.Where("phone = ?", req.Phone).First(&user).Error
		}
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(
			strconv.FormatUint(uint64(user.ID), 10),
			user.Email,
			time.Hour*24,
			&cfg.JWT,
		)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		res := types.LoginResponse{
			ID:    user.ID,
			Email: user.Email,
			Phone: user.Phone,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
