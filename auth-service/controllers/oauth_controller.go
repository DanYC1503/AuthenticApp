package controllers

import (
	"fmt"
	"log"
	"main/config"
	"main/middleware/encryption"
	repository "main/repository/auth-processing"
	"net/http"
	"net/url"

	"github.com/markbates/goth/gothic"
)

func GoogleCallback(w http.ResponseWriter, r *http.Request) {

	// Complete Google OAuth
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("OAuth failed in CompleteUserAuth: %v", err)
		http.Error(w, "OAuth failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	log.Printf("OAuth successful for user: %s (%s)", user.Name, user.Email)
	db := config.ConnectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database transaction failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Upsert user info (insert if new, update if existing)
	err = repository.UpsertOAuthUser(tx, user)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed: "+err.Error(), http.StatusInternalServerError)
		return
	}


	dbAfter := config.ConnectDB()
	defer dbAfter.Close()

	txAfter, err := dbAfter.Begin()
	if err != nil {
		http.Error(w, "Database transaction failed", http.StatusInternalServerError)
		return
	}
	defer txAfter.Rollback()

	// Generate session token
	sessionToken, expireDate, err := encryption.GenerateToken(user.Email, "session")
	if err != nil {
		http.Error(w, "Could not generate session token", http.StatusInternalServerError)
		return
	}

	// Check user type
	userType, err := repository.ReturnUserType(txAfter, user.Email)
	if err != nil {
		userType = "client" // fallback default
	}

	// Check user active status
	active, err := repository.ReturnUserStatus(txAfter, user.Email)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !active {
		log.Printf("User account is not active: %s", user.Email)
		http.Error(w, "User account is not active", http.StatusForbidden)
		return
	}

	// Ccleanly end the transaction
	if err := txAfter.Commit(); err != nil {
		log.Printf("Warning: failed to commit read-only tx: %v", err)
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
	frontendURL := fmt.Sprintf(
		"http://localhost:4200/ClientDashboard?username=%s&email=%s&userType=%s",
		url.QueryEscape(user.Email),
		url.QueryEscape(user.Email),
		url.QueryEscape(userType),
	)
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}
