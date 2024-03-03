// FOR LOCAL DEVELOPMENT
package main

import (
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/server"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize environment variables
	godotenv.Load(".env")

	// Initialize Auth
	auth.InitAuth()

	// Initialize Database
	db.InitDB()

	// Register routes
	app := server.RegisterRoutes()
	app.Run(":8080")
}
