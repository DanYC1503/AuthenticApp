package config

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	sessionKey = "supersecret-session-key"
	IsProd     = false
)

func InitOAuth() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env not loaded, using system environment variables")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}

	// Configure session store
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteNoneMode
	store.Options.Secure = false

	gothic.Store = store

	// This is where you register OAuth providers
	goth.UseProviders(
		google.New(
			clientID,
			clientSecret,
			"http://localhost:9999/auth/google/callback",
			"email", "profile",
		),
	)
}
