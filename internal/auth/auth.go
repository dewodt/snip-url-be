package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

// Auth Constants
const (
	maxAge = 86400 * 7
)

func NewAuth() {
	// Environment
	isProduction := os.Getenv("ENVIRONMENT") == "production"
	SESSION_SECRET := os.Getenv("SESSION_SECRET")

	// Google
	googleOAuthClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleOAuthClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	googleOAuthCallbackURL := os.Getenv("BE_URL") + "/api/auth/google/callback"

	// Discord
	discordOAuthClientID := os.Getenv("DISCORD_OAUTH_CLIENT_ID")
	discordOAuthClientSecret := os.Getenv("DISCORD_OAUTH_CLIENT_SECRET")
	discordOAuthCallbackURL := os.Getenv("BE_URL") + "/api/auth/discord/callback"

	// Github
	githubOAuthClientID := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	githubOAuthClientSecret := os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	githubOAuthCallbackURL := os.Getenv("BE_URL") + "/api/auth/github/callback"

	// Configure cookie
	store := sessions.NewCookieStore([]byte(SESSION_SECRET))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProduction // Local uses http (secure = false)

	// Goth config
	gothic.Store = store
	goth.UseProviders(
		google.New(googleOAuthClientID, googleOAuthClientSecret, googleOAuthCallbackURL),
		discord.New(discordOAuthClientID, discordOAuthClientSecret, discordOAuthCallbackURL),
		github.New(githubOAuthClientID, githubOAuthClientSecret, githubOAuthCallbackURL),
	)
}
