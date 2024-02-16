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
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func OAuthCallbackHandler(c *gin.Context) {
	// Get provider and set it in context (gothic uses context to get provider)
	provider := c.Param("provider")
	c.Request = gothic.GetContextWithProvider(c.Request, provider)

	// Handle Auth
	userAuth, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	// Check for errors
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is already registered
	var userDB models.User
	dbRes := db.DB.Where("email = ?", userAuth.Email).First(&userDB)
	// Check for errors
	if dbRes.Error != nil && !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// User is not registered
	if dbRes != nil && errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		// Create user
		userDB = models.User{
			Email: userAuth.Email,
			Name:  userAuth.Name,
		}
		dbRes = db.DB.Create(&userDB)
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
		"id":    userDB.ID,
		"email": userDB.Email,
		"name":  userDB.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSigned, err := jwtToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	// Save to cookie
	c.SetCookie("auth", jwtSigned, int(time.Hour*24), "/", os.Getenv("FE_URL"), false, true)

	// Redirect to frontend
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL"))
}
