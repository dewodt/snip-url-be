package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// Get claims from context
func GetClaimsFromContext(c *gin.Context) map[string]interface{} {
	// Check if user is authenticated from JWT
	jwtSigned, err := c.Cookie("auth")
	if err != nil {
		return nil
	}

	// Parse jwt
	JWT_SECRET := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtSigned, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil
	}

	return claims
}
