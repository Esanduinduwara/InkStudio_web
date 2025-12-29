package database

import (
	"database/sql"
	"fmt"
	"inkstudio-backend/internal/models"
)

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, email, username, passwordHash string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, created_at, updated_at
	`, email, username, passwordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}
