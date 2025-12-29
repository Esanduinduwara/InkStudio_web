package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt (includes automatic salt generation)
func HashPassword(password string) (string, error) {
	// bcrypt automatically generates a unique salt for each password
	// Cost of 12 provides a good balance between security and performance
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword checks if a password matches the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
