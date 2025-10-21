package models

import (
	"time"
)

type UserLoggedIn struct {
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	DateOfBirth   DateOnly  `json:"date_of_birth"`
	Address       string    `json:"address"`
	CreateDate    DateOnly  `json:"create_date"`
	Username      string    `json:"username"`
	IsVerified    bool      `json:"is_verified"`
	LastLogin     time.Time `json:"last_login"`
	AccountStatus string    `json:"account_status"`
}
type UserUpdate struct {
	Username          string   `json:"username"`
	IDNumberEncrypted string   `json:"id"`
	FullName          string   `json:"full_name"`
	Email             string   `json:"email"`
	PhoneNumber       string   `json:"phone_number"`
	DateOfBirth       DateOnly `json:"date_of_birth"`
	Address           string   `json:"address"`
	UsernameNew       string   `json:"username_new"`
}
type UserCredentialUpdate struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	PasswordNew string `json:"password_new"`
}
type UserRequestInfo struct {
	Username string `json:"username"`
}

type AdminUserEnableDisable struct {
	Username       string `json:"username"`
	ClientUsername string `json:"client_username"`
}
