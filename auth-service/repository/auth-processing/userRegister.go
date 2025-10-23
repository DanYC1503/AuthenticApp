package repository

import (
	"database/sql"
	"fmt"
	"main/middleware/encryption"
	"main/models"
	"time"

	_ "github.com/lib/pq"
)

func InsertUser(tx *sql.Tx, userCreateClient models.UserCreateClient, createDate time.Time) error {
	fmt.Println("Inserting User")

	query := `INSERT INTO users (
		id_number_encrypted, 
		full_name, 
		email, 
		phone_number, 
		date_of_birth, 
		address, 
		create_date, 
		username, 
		password_hash,
		salt, 
		is_verified, 
		last_login, 
		account_status,
		oauth_provider,
		oauth_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	// Generate salt and hash password
	salt, err := encryption.GenerateSalt(16)
	if err != nil {
		fmt.Printf("Salt generation error: %v\n", err)
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	passwordHash := encryption.HashPassword(userCreateClient.Password, salt)
	fmt.Printf("Password hashed successfully\n")

	// Encrypt ID number
	idNumberEncrypted, err := encryption.EncryptIDNumber(userCreateClient.IDNumber)
	if err != nil {
		return fmt.Errorf("failed to encrypt ID number: %w", err)
	}
	fmt.Printf("ID number encrypted successfully\n")

	// Handle optional fields
	var lastLogin sql.NullTime
	var dob sql.NullTime
	var oauthProvider sql.NullString
	var oauthID sql.NullString

	if !userCreateClient.DateOfBirth.ToTime().IsZero() {
		dob = sql.NullTime{Time: userCreateClient.DateOfBirth.ToTime(), Valid: true}
	}

	// Execute insert
	_, err = tx.Exec(query,
		idNumberEncrypted,
		userCreateClient.FullName,
		userCreateClient.Email,
		userCreateClient.PhoneNumber,
		dob,
		userCreateClient.Address,
		createDate,
		userCreateClient.Username,
		passwordHash,
		salt,
		false,
		lastLogin,
		"active",
		oauthProvider,
		oauthID,
	)

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
