package auth

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expires_at"`
}

// Get claims from context
func GetSessionFromContext(c *gin.Context) *Session {
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

	// Bind claims to session
	session := Session{}
	jsonString, err := json.Marshal(claims)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(jsonString, &session)
	if err != nil {
		return nil
	}

	return &session
}
