package backend
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type App struct {
	DB        *sql.DB
	JWTSecret []byte
}

func main() {
	// Load configuration from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default local database URL
		dbURL = "postgres://postgres:postgres@localhost:5432/inkstudio?sslmode=disable"
		log.Println("DATABASE_URL not set, using default:", dbURL)
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
		log.Println("JWT_SECRET not set, using default (change in production!)")
	}

	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize database tables
	if err := initDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app := &App{
		DB:        db,
		JWTSecret: []byte(jwtSecret),
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", app.RegisterHandler)
	mux.HandleFunc("/api/login", app.LoginHandler)
	mux.HandleFunc("/api/profile", app.AuthMiddleware(app.ProfileHandler))
	mux.HandleFunc("/api/health", app.HealthHandler)

	// CORS middleware wrapper
	handler := corsMiddleware(loggingMiddleware(mux))

	// Create HTTP server
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("📝 Endpoints:")
	fmt.Println("   POST /api/register - Register a new user")
	fmt.Println("   POST /api/login    - Login with credentials")
	fmt.Println("   GET  /api/profile  - Get user profile (requires auth)")
	fmt.Println("   GET  /api/health   - Health check")
	
	log.Fatal(srv.ListenAndServe())
}

func initDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	
	log.Println("✓ Database tables initialized successfully")
	return nil
}
