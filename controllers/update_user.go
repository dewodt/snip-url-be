package controllers

import (
	"net/http"
	"os"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UpdateUserSchema struct {
	Name   string  `form:"name" binding:"required"`
	Avatar *string `form:"avatar"`
}

func UpdateUserHandler(c *gin.Context) {
	// Validate & bind form data
	formData := UpdateUserSchema{}
	err := c.ShouldBind(&formData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Update user settings
	err = db.DB.Model(&models.User{}).Select("name", "avatar").Where("id = ?", session.ID).Updates(models.User{
		Name:   formData.Name,
		Avatar: formData.Avatar,
	}).Error

	// Check for errors
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user settings"})
		return
	}

	// Update session (jwt token)
	JWT_SECRET := os.Getenv("JWT_SECRET")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         session.ID,
		"email":      session.Email,
		"name":       formData.Name,
		"avatar":     formData.Avatar,
		"expires_at": time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSigned, err := jwtToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	// Save to cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("snip-url-auth", jwtSigned, 24*3600, "/", os.Getenv("PARENT_DOMAIN"), true, true)

	// Success
	c.JSON(http.StatusOK, gin.H{"message": "User settings updated"})
}
