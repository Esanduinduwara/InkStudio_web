package database

import (
	"database/sql"
	"fmt"
	"inkstudio-backend/internal/models"
)

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, email, username, passwordHash string) (*models.User, error) {
	// Insert user into database
	result, err := db.Exec(`
		INSERT INTO users (email, username, password_hash)
		VALUES (?, ?, ?)
	`, email, username, passwordHash)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the last inserted ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Fetch the created user
	var user models.User
	err = db.QueryRow(`
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
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
		WHERE id = ?
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
