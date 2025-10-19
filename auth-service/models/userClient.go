package models

import (
	"time"
)

type UserClient struct {
	IDNumberEncrypted string    `json:"id_number_encrypted"` // BYTEA
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PhoneNumber       string    `json:"phone_number,omitempty"`
	DateOfBirth       time.Time `json:"date_of_birth,omitempty"`
	Address           string    `json:"address,omitempty"`
	CreateDate        time.Time `json:"create_date"`
	Username          string    `json:"username"`
	PasswordHash      []byte    `json:"password_hash"` // BYTEA
	Salt              []byte    `json:"salt"`          // BYTEA
	IsVerified        bool      `json:"is_verified"`
	AccountStatus     string    `json:"account_status"`
}

type UserCreateClient struct {
	IDNumber    string   `json:"id_number"` // Will be encrypted before storing
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	DateOfBirth DateOnly `json:"date_of_birth,omitempty"`
	Address     string   `json:"address,omitempty"`
	Username    string   `json:"username"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserLoggedIn struct {
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Address       string    `json:"address"`
	CreateDate    time.Time `json:"create_date"`
	Username      string    `json:"username"`
	IsVerified    bool      `json:"is_verified"`
	LastLogin     time.Time `json:"last_login"`
	AccountStatus string    `json:"account_status"`
}
