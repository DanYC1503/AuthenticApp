package repository

import (
	"database/sql"
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
        account_status, 
        role
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	var DateNow = createDate
	var IsVerified = false
	var AccountStatus = "active"
	var Role = "user"
	salt, erro := encryption.GenerateSalt(16)
	if erro != nil {
		log.Fatal("Failed to generate salt:", erro)
	}
	user_password_hash := encryption.HashPassword(userCreateClient.Password, salt)
	IDNumberEncrypted := encryption.HashId_Number(string(userCreateClient.IDNumber), salt)
	var pk string
	err := tx.QueryRow(query,
		IDNumberEncrypted,
		userCreateClient.FullName,
		userCreateClient.Email,
		userCreateClient.PhoneNumber,
		userCreateClient.DateOfBirth,
		userCreateClient.Address,
		DateNow,
		userCreateClient.Username,
		[]byte(user_password_hash),
		salt,
		IsVerified,
		DateNow,
		AccountStatus,
		Role).Scan(&pk)
	if err != nil {
		log.Fatal("Failed to create User DB:", err)
	}

	return pk, err
}
