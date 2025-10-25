package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/config"
	"main/middleware"
	"main/middleware/encryption"
	"main/models"

	"net/http"
)

func RequestPasswordAuthToken(w http.ResponseWriter, r *http.Request, user models.UserPasswordRetrieval) {
	db := config.ConnectDB()
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT username FROM users WHERE email=$1", user.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		fmt.Println("SELECT error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resetToken, _, _ := encryption.GenerateToken(userID, "reset")
	if err != nil {
		http.Error(w, "Could not generate reset token", http.StatusInternalServerError)
		return
	}

	// Send email
	if err := middleware.SendPasswordRecoveryEmail(user.Email, resetToken); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the email exists, a recovery link has been sent.",
	})
}

func ResetPassword(w http.ResponseWriter, userInfo models.UserPasswordReset) error {
	db := config.ConnectDB()
	defer db.Close()

	// Debug: print the email received
	fmt.Println("ResetPassword called with email:", userInfo.Email)

	var salt []byte
	err := db.QueryRow(`SELECT salt FROM users WHERE email=$1`, userInfo.Email).Scan(&salt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return err
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return err
	}

	hashedPassword := encryption.HashPassword(userInfo.New_Password, salt)
	_, err = db.Exec(`UPDATE users SET password_hash=$1 WHERE email=$2`, hashedPassword, userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return err
	}

	fmt.Println("Password updated successfully for:", userInfo.Email)
	return nil
}
