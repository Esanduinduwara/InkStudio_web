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
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
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
		log.Printf("Failed to get user: %v", err)
		response.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		log.Printf("Invalid password for user %s", req.Email)
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

	// Clear password hash before sending response
	user.PasswordHash = ""

	response.JSON(w, models.AuthResponse{
		Token:   token,
		User:    *user,
		Message: "Login successful",
	}, http.StatusOK)
}

// GetMe returns the current user's information (protected route)
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user from database
	user, err := database.GetUserByID(h.DB, userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		response.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Clear password hash before sending response
	user.PasswordHash = ""

	response.JSON(w, user, http.StatusOK)
}
