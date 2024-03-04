package controllers

import (
	"errors"
	"net/http"
	"os"
	"snip-url-be/db"
	"snip-url-be/emails"
	"snip-url-be/models"
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
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL")+"/auth/error?"+err.Error())
		return
	}

	// Check if user is already registered
	var userDB models.User
	err = db.DB.Where("email = ?", userAuth.Email).First(&userDB).Error
	// Check for errors
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL")+"/auth/error?failed to get user")
		return
	}
	// User is not registered
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create user
		userDB = models.User{
			Email: userAuth.Email,
			Name:  userAuth.Name,
		}
		err = db.DB.Create(&userDB).Error
		// Check for errors
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL")+"/auth/error?failed to create user")
			return
		}

		// Send welcome email
		_, err := emails.SendWelcomeEmail(userDB.Email)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL")+"/auth/error?failed to send welcome email")
			return
		}
	}

	// Sign in user
	// Create jwt
	JWT_SECRET := os.Getenv("JWT_SECRET")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         userDB.ID,
		"email":      userDB.Email,
		"name":       userDB.Name,
		"avatar":     userDB.Avatar,
		"expires_at": time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSigned, err := jwtToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL")+"/auth/error?failed to sign token")
		return
	}

	// Save to cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("snip-url-auth", jwtSigned, 24*3600, "/", os.Getenv("PARENT_DOMAIN"), true, true)

	// Redirect to frontend
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FE_URL"))
}
