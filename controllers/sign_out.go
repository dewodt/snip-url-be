package controllers

import (
	"net/http"

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
	c.SetCookie("auth", "", -1, "/", "", false, true)

	// Success
	c.JSON(http.StatusOK, gin.H{"message": "Signed out"})
}
