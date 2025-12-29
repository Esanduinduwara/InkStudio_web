package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"inkstudio-backend/config"
	"inkstudio-backend/internal/database"
	"inkstudio-backend/internal/handlers"
	"inkstudio-backend/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Initialize handlers
	handler := handlers.NewHandler(db, cfg.JWTSecret)

	// Setup routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", handler.HealthCheck)

	// Auth routes
	mux.HandleFunc("/api/auth/register", handler.Register)
	mux.HandleFunc("/api/auth/login", handler.Login)

	// Protected route example
	mux.HandleFunc("/api/auth/me", middleware.AuthMiddleware(handler.GetMe, cfg.JWTSecret))

	// Apply middleware
	corsHandler := middleware.CORS(mux)
	loggedHandler := middleware.Logging(corsHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, loggedHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
