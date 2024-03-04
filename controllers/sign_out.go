package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func SignOutHandler(c *gin.Context) {
	// Sign out user
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign out"})
		return
	}

	// Clear cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("snip-url-auth", "", -1, "/", os.Getenv("PARENT_DOMAIN"), true, true)

	// Success
	c.JSON(http.StatusOK, gin.H{"message": "Signed out"})
}
