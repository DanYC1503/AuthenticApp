package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	_ "github.com/lib/pq"
)

func GenerateSalt(length int) ([]byte, error) {
	saltBytes := make([]byte, length)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return saltBytes, nil
}
func HashPassword(password string, salt []byte) string {
	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hash.Write([]byte(password))
	hashedBytes := hash.Sum(nil)

	return hex.EncodeToString(hashedBytes)
}
func HashId_Number(id_number string, salt []byte) string {
	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(id_number))
	hash.Write([]byte(id_number))
	hashedBytes := hash.Sum(nil)

	return hex.EncodeToString(hashedBytes)
}
