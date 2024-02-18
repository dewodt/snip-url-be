package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate a secure token (for email)
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	// Error handling
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
