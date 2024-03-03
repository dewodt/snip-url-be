package controllers

import (
	"errors"
	"net/http"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/emails"
	"snip-url-be/models"

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
	err := c.ShouldBind(&formData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Get data
	emailData := formData.Email
	nameData := formData.Name

	// Check if user is already registered
	var user models.User
	err = db.DB.Where("email = ?", emailData).First(&user).Error
	// User is already registered
	if user.Email == emailData {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already registered", "field": "email"})
		return
	}
	// Check for other errors
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
		return
	}

	// Generate token
	token, err := auth.GenerateSecureToken(64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save email & token in DB
	verification := models.Verification{
		Email: emailData,
		Name:  &nameData,
		Token: token,
	}
	err = db.DB.Create(&verification).Error
	// Failed to save token
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Send email
	_, err = emails.SendSignInEmail(emailData, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}
