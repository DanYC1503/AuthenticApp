package encryption

import (
	"errors"
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

	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") {
		return true
	}

	return false
}
func ValidateToken(tokenString string, expectedType string) (*models.Claims, bool) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}

	// Check that the token type matches what we expect
	if claims.TokenType != expectedType {
		return nil, false
	}

	return claims, true
}
