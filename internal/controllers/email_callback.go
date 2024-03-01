package controllers

import (
	"errors"
	"net/http"
	"os"
	"snip-url-be/internal/db"
	"snip-url-be/internal/emails"
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
	token := c.Query("token")
	email := c.Query("email")

	// Check if token and email is valid string
	if token == "" || email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token or email"})
		return
	}

	// Validate token
	var verification models.Verification
	err := db.DB.Where("token = ? AND email = ? AND expires_at > ?", token, email, time.Now()).First(&verification).Error
	// Invalid token
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}
	// Check for other errors
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	// Check if user is already registered
	var user models.User
	err = db.DB.Where("email = ?", verification.Email).First(&user).Error
	// Check query errors
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	// Check if user is not found
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Create user
		user = models.User{
			Email: verification.Email,
			Name:  *verification.Name,
		}
		err = db.DB.Create(&user).Error
		// Check for errors
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Send welcome email
		_, err = emails.SendWelcomeEmail(user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
			return
		}
	}

	// Sign in user
	// Create jwt
	JWT_SECRET := os.Getenv("JWT_SECRET")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"avatar":     user.Avatar,
		"expires_at": time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSigned, err := jwtToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	// Save to cookie
	c.SetCookie("auth", jwtSigned, 24*3600, "/", os.Getenv("FE_URL"), false, true)

	// Return success
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL"))
}
