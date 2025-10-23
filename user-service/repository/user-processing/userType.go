package repository

import (
	"database/sql"
	"fmt"
)

func VerifyUserType(db *sql.DB, username string) (bool, error) {
	query := `SELECT user_type FROM users WHERE username = $1`

	var user_type string
	err := db.QueryRow(query, username).Scan(&user_type)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, err
	}

	// Return true if account is active
	if user_type == "admin" {
		return true, nil
	}

	// Return false if account is disabled or any other status
	return false, nil
}
func ReturnUserType(tx *sql.Tx, username string) (string, error) {
	query := `SELECT user_type FROM users WHERE username = $1`

	var userType string
	err := tx.QueryRow(query, username).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to query user_type: %w", err)
	}

	return userType, nil
}
