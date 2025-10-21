package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"

	_ "github.com/lib/pq"
)

func DisableUser(w http.ResponseWriter, user models.AdminUserEnableDisable, tx *sql.Tx) error {
	// Verify that the requester is an admin
	if err := VerifyAdminStatus(tx, user.Username); err != nil {
		http.Error(w, "Unauthorized - only admin can disable users", http.StatusForbidden)
		return err
	}

	// Disable the target user
	query := `UPDATE users 
	          SET account_status = 'disabled'
	          WHERE username = $1 
	          RETURNING username`

	var disabledUsername string
	err := tx.QueryRow(query, user.ClientUsername).Scan(&disabledUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return fmt.Errorf("user %s not found", user.ClientUsername)
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("failed to disable user: %w", err)
	}

	// Respond with confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"disabled_user": disabledUsername,
	})

	return nil
}
func EnableUser(w http.ResponseWriter, user models.AdminUserEnableDisable, tx *sql.Tx) error {
	// Verify that the requester is an admin
	if err := VerifyAdminStatus(tx, user.Username); err != nil {
		http.Error(w, "Unauthorized - only admin can enable users", http.StatusForbidden)
		return err
	}

	// Disable the target user
	query := `UPDATE users 
	          SET account_status = 'active'
	          WHERE username = $1 
	          RETURNING username`

	var enabledUsername string
	err := tx.QueryRow(query, user.ClientUsername).Scan(&enabledUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return fmt.Errorf("user %s not found", user.ClientUsername)
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("failed to disable user: %w", err)
	}

	// Respond with confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"enabled_user": enabledUsername,
	})

	return nil
}

func VerifyAdminStatus(tx *sql.Tx, username string) error {
	var userType string
	err := tx.QueryRow(`SELECT user_type FROM users WHERE username=$1`, username).Scan(&userType)
	if err != nil {
		return err
	}
	if userType != "admin" {
		return fmt.Errorf("user %s is not an admin", username)
	}
	return nil
}
