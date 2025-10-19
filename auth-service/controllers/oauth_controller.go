package controllers

import (
	"encoding/json"
	"log"
	"main/config"
	"main/middleware/encryption"
	repository "main/repository/auth-processing"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func GoogleCallback(w http.ResponseWriter, r *http.Request) {

	// Complete Google OAuth
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "OAuth failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	db := config.ConnectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database transaction failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Upsert user info (insert if new, update if existing)
	userID, err := repository.UpsertOAuthUser(tx, user)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("OAuth user processed, ID:", userID)
	err = repository.UpdateAuthMethodLastUsed(db, userID, "oauth")
	if err != nil {
		log.Println("Warning: failed to update auth_methods.last_used:", err)
	}

	// Generate session token
	sessionToken, expireDate, err := encryption.GenerateSessionToken(user.Email)
	if err != nil {
		http.Error(w, "Could not generate session token", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expireDate,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	// Instead of redirect, respond like your normal login
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_token": sessionToken,
		"expires":       expireDate.Unix(),
		"user_email":    user.Email,
	})
}
