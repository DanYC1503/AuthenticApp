package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"main/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT for the given username
func GenerateSessionToken(username string) (string, time.Time, error) {
	expSeconds, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION")) // expiration in seconds
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if err != nil || expSeconds <= 0 {
		expSeconds = 3600 // fallback: 1 hour
	}
	expirationTime := time.Now().Add(time.Duration(expSeconds) * time.Second)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "session",
		},
	}

	// sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func GenerateActionJWT(username, action string) (string, time.Time, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	var subject string
	var duration time.Duration

	switch action {
	case "delete":
		subject = "deleteAuth"
		duration = 5 * time.Minute
	case "update":
		subject = "updateAuth"
		duration = 10 * time.Minute
	case "auth":
		subject = "passwordAuth"
		duration = 15 * time.Minute
	default:
		return "", time.Time{}, fmt.Errorf("invalid action type: %s", action)
	}

	expirationTime := time.Now().Add(duration)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}
func GenerateManualCSRFToken() string {
	// Generate a random 32-byte token
	token := make([]byte, 32)
	rand.Read(token)
	return base64.StdEncoding.EncodeToString(token)
}
