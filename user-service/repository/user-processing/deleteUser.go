package repository

import (
	"database/sql"
	"fmt"
	"main/models"
	"net/http"

	_ "github.com/lib/pq"
)

func DeleteUser(w http.ResponseWriter, user models.UserRequestInfo, tx *sql.Tx) error {
	query := `DELETE FROM users WHERE username = $1 RETURNING username`

	var deletedUsername string
	err := tx.QueryRow(query, user.Username).Scan(&deletedUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return fmt.Errorf("user %s not found", user.Username)
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
