package server

import (
	"net/http"
	"snip-url-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			// Email provider
			auth.POST("/email/sign-in", controllers.EmailSignInProviderHandler)
			auth.POST("/email/sign-up", controllers.EmailSignUpProviderHandler)
			auth.GET("/email/callback/:token", controllers.EmailCallbackHandler)

			// OAuth provider
			auth.GET("/:provider", controllers.OAuthProviderHandler)
			auth.GET("/:provider/callback", controllers.OAuthCallbackHandler)

			// Check session
			auth.GET("/session", controllers.SessionHandler)

			// Sign out
			auth.GET("/sign-out", controllers.SignOutHandler)
		}

		// Upload file
		api.POST("/upload-avatar", controllers.UploadAvatarHandler)

		// Snip
	}

	return r
}
