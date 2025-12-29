package handlers

import "database/sql"

// Handler holds dependencies for HTTP handlers
type Handler struct {
	DB        *sql.DB
	JWTSecret []byte
}

// NewHandler creates a new Handler instance
func NewHandler(db *sql.DB, jwtSecret string) *Handler {
	return &Handler{
		DB:        db,
		JWTSecret: []byte(jwtSecret),
	}
}
