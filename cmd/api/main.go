package main

import (
	"fmt"
	"snip-url-be/internal/auth"
	"snip-url-be/internal/db"
	"snip-url-be/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize Env
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Error loading .env file")
	}

	// Initialize Auth
	auth.InitAuth()

	// Initialize Database
	db.InitDB()

	// Initialize Server
	server := server.NewServer()
	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(fmt.Sprintf("cannot start server: %s", errServer))
	}
}
