package server

import (
	"net/http"
	"os"
	"snip-url-be/internal/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FE_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

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
