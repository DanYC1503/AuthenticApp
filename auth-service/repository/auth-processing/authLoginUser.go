package repository

import (
	"database/sql"
	"fmt"
	"main/middleware/encryption"
	"net/http"

	_ "github.com/lib/pq"
)

func AuthenticateUserCredentials(db *sql.DB, username, password string) (bool, string, error) {
	query := `SELECT username, password_hash, salt FROM users WHERE username = $1`

	var dbUsername string
	var passwordHash, salt []byte

	err := db.QueryRow(query, username).Scan(&dbUsername, &passwordHash, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil // user not found
		}
		return false, "", err
	}

	ok := encryption.VerifyPasswordBySaltAndHash(salt, passwordHash, password)
	if !ok {
		return false, "", nil
	}
	_, err = db.Exec(`UPDATE users SET last_login = NOW() WHERE username = $1`, username)
	if err != nil {
		return false, "", fmt.Errorf("failed to update last_login: %w", err)
	}

	_, err = db.Exec(`SELECT update_auth_method_last_login($1)`, username)
	if err != nil {
		return false, "", fmt.Errorf("failed to update auth_methods.last_used: %w", err)
	}

	return true, dbUsername, nil
}

func SessionTokenVerification(w http.ResponseWriter, r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Missing session token", http.StatusUnauthorized)
		return false, ""
	}

	//  Validate session token and get username
	username, err := encryption.ValidateSessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid or expired session token", http.StatusUnauthorized)
		return false, ""
	}
	fmt.Printf("Session Token Verified")
	return true, username
}
