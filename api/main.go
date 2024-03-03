// FOR PRODUCTION
package api

import (
	"net/http"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/server"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	// Initialize Auth
	auth.InitAuth()

	// Initialize Database
	db.InitDB()

	// Register routes
	app = server.RegisterRoutes()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize Entrypoint
	app.ServeHTTP(w, r)
}
