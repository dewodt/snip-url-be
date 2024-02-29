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

// Form data schema for sign up with email
type SignUpEmailSchema struct {
	Email string `form:"email" binding:"required,email"`
	Name  string `form:"name" binding:"required"`
}

// Route: /auth/email/sign-up
func EmailSignUpProviderHandler(c *gin.Context) {
	// Validate form data
	var formData SignUpEmailSchema
	bindErr := c.ShouldBind(&formData)
	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Get data
	emailData := formData.Email
	nameData := formData.Name

	// Check if user is already registered
	var user models.User
	dbRes := db.DB.Where("email = ?", emailData).First(&user)
	// User is already registered
	if dbRes.RowsAffected > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already registered", "field": "email"})
		return
	}

	// Check for other errors
	if dbRes.Error != nil && !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
		return
	}

	// Generate token
	token, tokenErr := auth.GenerateSecureToken(64)
	if tokenErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save token in DB
	verification := models.Verification{
		Email: emailData,
		Name:  &nameData,
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
