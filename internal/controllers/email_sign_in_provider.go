package controllers

import (
	"errors"
	"net/http"
	"snip-url-be/internal/auth"
	"snip-url-be/internal/db"
	"snip-url-be/internal/emails"
	"snip-url-be/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Form data schema for sign in with email
type SignInEmailSchema struct {
	Email string `form:"email" binding:"required,email"`
}

// Route: /auth/email/sign-in
func EmailSignInProviderHandler(c *gin.Context) {
	// Validate form data
	var formData SignInEmailSchema
	bindErr := c.ShouldBind(&formData)
	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}
	// Get data
	emailData := formData.Email

	// Check if user is already registered
	var user models.User
	dbRes := db.DB.Where("email = ?", emailData).First(&user)
	// User is not registered
	if errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not registered", "field": "email"})
		return
	}
	// Check for other errors
	if dbRes.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
		return
	}

	// Generate token
	token, tokenErr := auth.GenerateSecureToken(64)
	if tokenErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save token & email in DB
	verification := models.Verification{
		Email: emailData,
		Token: token,
	}
	dbRes = db.DB.Create(&verification)
	// Failed to save token
	if dbRes.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Send email
	_, emailErr := emails.SendSignInEmail(emailData, token)
	if emailErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}
