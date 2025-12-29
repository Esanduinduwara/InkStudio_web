package handlers

import "database/sql"

type Handler struct {
	DB        *sql.DB
	JWTSecret string
}

func NewHandler(db *sql.DB, jwtSecret string) *Handler {
	return &Handler{
		DB:        db,
		JWTSecret: jwtSecret,
	}
}
