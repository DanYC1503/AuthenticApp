package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"

	_ "github.com/lib/pq"
)

func RetrieveUserLoginInfo(w http.ResponseWriter, user models.UserRequestInfo, tx *sql.Tx) error {
	var endUser models.UserLoggedIn

	queryUser := `SELECT account_status, address, create_date, date_of_birth, email, full_name, is_verified, last_login, phone_number, username 
              FROM users WHERE username = $1`

	row := tx.QueryRow(queryUser, user.Username)
	err := row.Scan(
		&endUser.AccountStatus,
		&endUser.Address,
		&endUser.CreateDate,
		&endUser.DateOfBirth,
		&endUser.Email,
		&endUser.FullName,
		&endUser.IsVerified,
		&endUser.LastLogin,
		&endUser.PhoneNumber,
		&endUser.Username,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return fmt.Errorf("user %s not found", user.Username)
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("failed to scan user row: %w", err)
	}

	// Encode and send the user info
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"user": endUser,
	})

}
func RetrieveUserUsername(w http.ResponseWriter, tx *sql.Tx) error {
	query := `SELECT username FROM users WHERE user_type = 'client'`

	rows, err := tx.Query(query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []string

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return fmt.Errorf("failed to scan username: %w", err)
		}
		users = append(users, username)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(users) == 0 {
		http.Error(w, "No clients found", http.StatusNotFound)
		return nil
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"clients": users,
	})
}

