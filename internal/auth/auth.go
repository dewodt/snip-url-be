package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
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
	googleOAuthCallbackURL := os.Getenv("BE_URL") + "/auth/google/callback"

	// Discord
	discordOAuthClientID := os.Getenv("DISCORD_OAUTH_CLIENT_ID")
	discordOAuthClientSecret := os.Getenv("DISCORD_OAUTH_CLIENT_SECRET")
	discordOAuthCallbackURL := os.Getenv("BE_URL") + "/auth/discord/callback"

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
	)
}
