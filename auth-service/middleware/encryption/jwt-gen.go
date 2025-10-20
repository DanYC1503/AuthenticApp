package encryption

import (
	"main/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT for the given username
func GenerateToken(username, tokenType string) (string, time.Time, error) {
	expSeconds, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION")) // expiration in seconds
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	if err != nil || expSeconds <= 0 {
		expSeconds = 1800 // fallback: 1 hour
	}
	expirationTime := time.Now().Add(time.Duration(expSeconds) * time.Second)

	claims := &models.Claims{
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   tokenType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, expirationTime, err
}
