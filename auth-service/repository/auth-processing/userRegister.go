package repository

import (
	"database/sql"
	"fmt"
	"log"
	"main/middleware/encryption"
	"main/models"

	"time"

	_ "github.com/lib/pq"
)

func InsertUser(tx *sql.Tx, userCreateClient models.UserCreateClient, createDate time.Time) (string, error) {
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
    account_status
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`

	salt, err := encryption.GenerateSalt(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	passwordHash := encryption.HashPassword(userCreateClient.Password, salt)
	idNumberEncrypted, err := encryption.EncryptIDNumber(userCreateClient.IDNumber)
	if err != nil {
		log.Fatal(err)
	}
	var lastLogin sql.NullTime // nil at creation
	var dob sql.NullTime
	if !userCreateClient.DateOfBirth.ToTime().IsZero() {
		dob = sql.NullTime{Time: userCreateClient.DateOfBirth.ToTime(), Valid: true}
	} else {
		dob = sql.NullTime{Valid: false}
	}

	var pk string
	err = tx.QueryRow(query,
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
		false,     // is_verified
		lastLogin, // last_login
		"active",  // account_status
	).Scan(&pk)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return pk, nil
}
