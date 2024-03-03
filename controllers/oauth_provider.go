package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OAuthProviderHandler(c *gin.Context) {
	// Get provider and set it in context (gothic uses context to get provider)
	provider := c.Param("provider")
	c.Request = gothic.GetContextWithProvider(c.Request, provider)

	// try to get the user without re-authenticating
	_, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	// Begin auth
	if err != nil {
		gothic.BeginAuthHandler(c.Writer, c.Request)
		return
	}
}
