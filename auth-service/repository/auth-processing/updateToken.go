package repository

import (
	"main/config"
	"main/middleware/encryption"
	"main/models"
	"net/http"
)

func RequestUpdateAuthToken(w http.ResponseWriter, r *http.Request, user models.UserLogin) string {
	db := config.ConnectDB()
	defer db.Close()

	// Issue short-lived token for update
	token, _, err := encryption.GenerateActionJWT(user.Username, "update")
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return ""
	}

	return token
}
