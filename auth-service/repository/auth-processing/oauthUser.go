package repository

import (
	"database/sql"
	"fmt"

	"github.com/markbates/goth"
)

func UpsertOAuthUser(tx *sql.Tx, googleUser goth.User) (string, error) {
	query := `
	INSERT INTO users (full_name, email, username, is_verified, create_date, account_status)
	VALUES ($1, $2, $3, true, NOW(), 'active')
	ON CONFLICT (email)
	DO UPDATE SET last_login = NOW()
	RETURNING id
	`

	var userID string
	err := tx.QueryRow(
		query,
		googleUser.Name,
		googleUser.Email,
		googleUser.NickName, // fallback for username
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("failed to upsert oauth user: %w", err)
	}

	return userID, nil
}

func UpdateAuthMethodLastUsed(db *sql.DB, userID string, methodType string) error {
	query := `
		UPDATE auth_methods
		SET last_used = NOW()
		WHERE user_id = $1 AND method_type = $2
	`

	res, err := db.Exec(query, userID, methodType)
	if err != nil {
		return fmt.Errorf("failed to update last_used: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no auth_method found for user_id=%s and method_type=%s", userID, methodType)
	}

	return nil
}
