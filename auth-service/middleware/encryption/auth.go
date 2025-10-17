package encryption

import (
	"errors"
	"fmt"
	"main/models"

	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

func VerifyPasswordBySaltAndHash(salt []byte, passHash []byte, inputPassword string) bool {
	var passInput = HashPassword(inputPassword, salt)
	if passInput == string(passHash) {
		return true
	}
	return false
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var pqErr *pq.Error
	if ok := errors.As(err, &pqErr); ok {

		switch string(pqErr.Code) {
		case "40001": // serialization_failure - transaction retry safe
			return true
		case "40P01": // deadlock_detected
			return true
		}
	}

	// For other error types, do custom string matching
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") {
		return true
	}

	return false
}
func ValidateToken(tokenString string, expectedSubject string) bool {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	return claims.Subject == expectedSubject
}
func ValidateSessionToken(tokenString string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired token")
	}

	if claims.Subject != "session" {
		return "", fmt.Errorf("invalid token subject")
	}

	return claims.Username, nil
}

func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
