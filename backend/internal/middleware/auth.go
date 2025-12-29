package middleware

import (
	"context"
	"net/http"
	"strings"

	"inkstudio-backend/internal/auth"
	"inkstudio-backend/pkg/response"
)

// AuthMiddleware validates JWT tokens and protects routes
func AuthMiddleware(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "username", claims.Username)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
