package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/inkstudio?parseTime=true&charset=utf8mb4"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:        getEnv("PORT", "8080"),
	}

	// Warn if using default values in production
	if cfg.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("⚠️  WARNING: Using default JWT_SECRET. Set a secure secret in production!")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
