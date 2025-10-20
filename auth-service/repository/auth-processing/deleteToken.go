package repository

import (
	"main/config"
	"main/middleware/encryption"
	"main/models"

	"net/http"
)

func RequestDeleteAuthToken(w http.ResponseWriter, r *http.Request, user models.UserLogin) string {
	db := config.ConnectDB()
	defer db.Close()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return ""
	}
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return ""
	}

	// Issue short-lived delete token
	token, _, err := encryption.GenerateToken(user.Username, "delete")
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return ""
	}

	return token
}
