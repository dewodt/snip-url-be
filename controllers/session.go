package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Route: GET /api/auth/session (check if user is authenticated)
func SessionHandler(c *gin.Context) {
	// Check if user is authenticated from JWT
	jwtSigned, err := c.Cookie("snip-url-auth")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse jwt
	JWT_SECRET := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtSigned, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JWT"})
		return
	}

	// Return claims
	c.JSON(http.StatusOK, claims)
}
