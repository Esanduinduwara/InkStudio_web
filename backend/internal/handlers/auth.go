package handlers
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"inkstudio-backend/internal/auth"
	"inkstudio-backend/internal/database"
	"inkstudio-backend/internal/models"
	"inkstudio-backend/pkg/response"
)

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		response.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		response.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Email, "@") {
		response.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash password with bcrypt (automatically includes salt)
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		response.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user in database
	user, err := database.CreateUser(h.DB, req.Email, req.Username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			response.Error(w, "Email or username already exists", http.StatusConflict)
			return
		}
		log.Printf("Failed to create user: %v", err)
		response.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Username, h.JWTSecret)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		response.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response.JSON(w, models.AuthResponse{
		Token:   token,
		User:    *user,
		Message: "User registered successfully",
	}, http.StatusCreated)
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		response.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := database.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		log.Printf("Database error: %v", err)
		response.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		response.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Username, h.JWTSecret)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		response.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response.JSON(w, models.AuthResponse{
		Token:   token,
		User:    *user,
		Message: "Login successful",
	}, http.StatusOK)
}

// Profile returns the authenticated user's profile
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by AuthMiddleware)
	userID := r.Context().Value("user_id").(int)

	// Fetch user from database
	user, err := database.GetUserByID(h.DB, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Database error: %v", err)
		response.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	response.JSON(w, user, http.StatusOK)
}
