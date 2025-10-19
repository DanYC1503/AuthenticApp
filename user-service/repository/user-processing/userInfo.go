package repository

import (
	"database/sql"
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": models.UserLoggedIn{},
		})
	}

	return nil
}
