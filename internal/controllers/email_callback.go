package controllers

import (
	"errors"
	"net/http"
	"os"
	"snip-url-be/internal/db"
	"snip-url-be/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Handle callback
// /auth/email/callback/:token
func EmailCallbackHandler(c *gin.Context) {
	// Get token from URL
	token := c.Param("token")

	// Validate token
	var verification models.Verification
	dbRes := db.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&verification)
	// Invalid token
	if errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}
	// Check for other errors
	if dbRes.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	// Check if user is already registered
	var user models.User
	dbRes = db.DB.Where("email = ?", verification.Email).First(&user)
	// Check query errors
	if dbRes.Error != nil && !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	// Check if user is not found
	if dbRes.Error != nil && errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		// Create user
		user = models.User{
			Email: verification.Email,
			Name:  *verification.Name,
		}
		dbRes = db.DB.Create(&user)
		// Check for errors
		if dbRes.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Sign in user
	// Create jwt
	JWT_SECRET := os.Getenv("JWT_SECRET")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSigned, err := jwtToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	// Save to cookie
	c.SetCookie("auth", jwtSigned, int(time.Hour*24), "/", os.Getenv("FE_URL"), false, true)

	// Return success
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL"))
}
