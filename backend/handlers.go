package backend
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// RegisterHandler handles user registration
func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		respondError(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		respondError(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Email, "@") {
		respondError(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash password with bcrypt (automatically includes salt)
	// Cost of 12 is a good balance between security and performance
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		respondError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	var user User
	err = app.DB.QueryRow(`
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, created_at, updated_at
	`, req.Email, req.Username, string(hashedPassword)).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondError(w, "Email or username already exists", http.StatusConflict)
			return
		}
		log.Printf("Failed to create user: %v", err)
		respondError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := app.generateToken(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	respondJSON(w, AuthResponse{
		Token:   token,
		User:    user,
		Message: "User registered successfully",
	}, http.StatusCreated)
}

// LoginHandler handles user login
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	var user User
	err := app.DB.QueryRow(`
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		respondError(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		respondError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := app.generateToken(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	respondJSON(w, AuthResponse{
		Token:   token,
		User:    user,
		Message: "Login successful",
	}, http.StatusOK)
}

// ProfileHandler returns the authenticated user's profile
func (app *App) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (set by AuthMiddleware)
	claims := r.Context().Value("user").(*Claims)

	// Fetch fresh user data from database
	var user User
	err := app.DB.QueryRow(`
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE id = $1
	`, claims.UserID).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondError(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		respondError(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	respondJSON(w, user, http.StatusOK)
}

// HealthHandler returns server health status
func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}, http.StatusOK)
}

// generateToken creates a JWT token for the user
func (app *App) generateToken(user User) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(app.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper functions for JSON responses
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, status int) {
	respondJSON(w, ErrorResponse{Error: message}, status)
}
