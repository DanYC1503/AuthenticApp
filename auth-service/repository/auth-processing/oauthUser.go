package repository

import (
	"database/sql"
	"fmt"

	"github.com/markbates/goth"
)

func UpsertOAuthUser(tx *sql.Tx, googleUser goth.User) error {
	userQuery := `
    INSERT INTO users (full_name, email, username, oauth_provider, oauth_id, is_verified, create_date, account_status)
    VALUES ($1, $2, $3, $4, $5, true, NOW(), 'active')
    ON CONFLICT (email)
    DO UPDATE SET 
        last_login = NOW(),
        oauth_provider = EXCLUDED.oauth_provider,
        oauth_id = EXCLUDED.oauth_id
    RETURNING id
    `

	// Use email as username to ensure consistency
	username := googleUser.Email

	var userID string
	err := tx.QueryRow(
		userQuery,
		googleUser.Name,
		googleUser.Email,
		username,
		"google",          // oauth_provider
		googleUser.UserID, // oauth_id
	).Scan(&userID)

	if err != nil {
		return fmt.Errorf("failed to upsert oauth user: %w", err)
	}

	return nil
}

func ReturnUserType(tx *sql.Tx, email string) (string, error) {
	query := `SELECT user_type FROM users WHERE email = $1`

	var userType string
	err := tx.QueryRow(query, email).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to query user_type: %w", err)
	}

	return userType, nil
}
func ReturnUserStatus(tx *sql.Tx, email string) (bool, error) {
	query := `SELECT account_status FROM users WHERE email = $1`

	var accountStatus string
	err := tx.QueryRow(query, email).Scan(&accountStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist yet: treat as inactive for now
			return false, nil
		}
		return false, fmt.Errorf("failed to query account_status: %w", err)
	}

	return accountStatus == "active", nil
}
