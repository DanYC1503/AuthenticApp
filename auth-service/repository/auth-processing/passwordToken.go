package repository

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/config"
	"main/middleware"
	"main/middleware/encryption"
	"main/models"
	"time"

	"github.com/google/uuid"

	"net/http"
)

func RequestPasswordAuthToken(w http.ResponseWriter, r *http.Request, user models.UserPasswordRetrieval) string {
	db := config.ConnectDB()
	defer db.Close()

	// Check if user exists and get ID
	var userID uuid.UUID
	err := db.QueryRow("SELECT id FROM users WHERE email=$1", user.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return ""
		}
		fmt.Println("SELECT error:", err)
		http.Error(w, "Database error request function", http.StatusInternalServerError)
		return ""
	} else {
		fmt.Println("User found:", userID)
	}

	// Generate secure random token hash
	token := uuid.New().String()             // or use crypto/rand for higher security
	tokenHash := encryption.HashToken(token) // hash for database storage

	// Save token into recovery_tokens
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	_, err = db.Exec(`
    INSERT INTO recovery_tokens (user_id, token_hash, expires_at, used, created_at)
    VALUES ($1, $2, $3, false, NOW() AT TIME ZONE 'UTC')
`, userID, tokenHash, expiresAt)

	// Send email - but don't fail the request if email fails
	if err := middleware.SendPasswordRecoveryEmail(user.Email, token); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)

	}

	return "If the email exists, a recovery link has been sent."
}
func VerifyRecoveryToken(w http.ResponseWriter, token string) {
	fmt.Printf("=== VerifyRecoveryToken called ===\n")

	db := config.ConnectDB()
	defer db.Close()

	tokenHash := encryption.HashToken(token)
	fmt.Printf("Hashed token (hex): %s\n", tokenHash)

	debugRecentTokens(db) // optional debug helper

	tokenHashBytes, err := decodeTokenHash(tokenHash)
	checkTokenStatus(db, tokenHashBytes)
	if err != nil {
		http.Error(w, "Invalid token format", http.StatusBadRequest)
		return
	}

	userID, expiresAt, used, err := getRecoveryTokenInfo(db, tokenHashBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := validateTokenStatus(expiresAt, used); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := markTokenUsed(db, tokenHashBytes); err != nil {
		http.Error(w, "Failed to update token status", http.StatusInternalServerError)
		return
	}

	resetToken, err := generateResetToken(userID)
	if err != nil {
		http.Error(w, "Could not generate reset token", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Success! Generated reset token for user %s\n", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"resetToken": resetToken,
		"message":    "Token verified successfully",
	})
}

// --- Helper functions ---

func debugRecentTokens(db *sql.DB) {
	rows, err := db.Query(`
        SELECT user_id, token_hash, used, expires_at, created_at 
        FROM recovery_tokens 
        ORDER BY created_at DESC LIMIT 5
    `)
	if err != nil {
		fmt.Printf("Error querying recovery_tokens: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("Recent tokens in recovery_tokens table:")
	for rows.Next() {
		var userID, tokenHash string
		var used bool
		var expiresAt, createdAt time.Time
		if err := rows.Scan(&userID, &tokenHash, &used, &expiresAt, &createdAt); err != nil {
			fmt.Printf("  Error scanning row: %v\n", err)
			continue
		}
		fmt.Printf("  ID: %s, UserID: %s, Hash: %s, Used: %t, Expires: %s\n",
			userID, tokenHash, used, expiresAt.Format(time.RFC3339))
	}
}

func decodeTokenHash(tokenHash string) ([]byte, error) {
	return hex.DecodeString(tokenHash)
}

func getRecoveryTokenInfo(db *sql.DB, tokenHash []byte) (uuid.UUID, time.Time, bool, error) {
	var userID uuid.UUID
	var expiresAt time.Time
	var used bool

	err := db.QueryRow(`
        SELECT user_id, expires_at, used 
        FROM recovery_tokens 
        WHERE token_hash=$1 AND used=false
    `, tokenHash).Scan(&userID, &expiresAt, &used)
	expiresAt = expiresAt.UTC()
	return userID, expiresAt, used, err
}

func validateTokenStatus(expiresAt time.Time, used bool) error {
	now := time.Now().UTC()
	if now.After(expiresAt) {
		return fmt.Errorf("Token expired")
	}
	if used {
		return fmt.Errorf("Token already used")
	}
	return nil
}

func markTokenUsed(db *sql.DB, tokenHash []byte) error {
	_, err := db.Exec(`UPDATE recovery_tokens SET used=true WHERE token_hash=$1`, tokenHash)
	return err
}

func generateResetToken(userID uuid.UUID) (string, error) {
	resetToken, _, err := encryption.GenerateToken(userID.String(), "reset")
	return resetToken, err
}
func checkTokenStatus(db *sql.DB, tokenHash []byte) {
	var userID string
	var used bool
	var expiresAt time.Time

	err := db.QueryRow(`
        SELECT user_id, used, expires_at 
        FROM recovery_tokens 
        WHERE token_hash=$1
    `, tokenHash).Scan(&userID, &used, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("  Token hash not found in database at all")
			return
		}
		fmt.Printf("  Error checking token status: %v\n", err)
		return
	}

	now := time.Now().UTC()
	fmt.Printf("  Token exists: UserID=%s, Used=%t, Expires=%s, Now(UTC)=%s\n",
		userID, used, expiresAt.Format(time.RFC3339), now.Format(time.RFC3339))
	fmt.Printf("  Is expired: %t\n", now.After(expiresAt))
}
