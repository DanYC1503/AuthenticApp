package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

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

func EncryptIDNumber(plainID string) (string, error) {
	key := []byte(os.Getenv("IDNUMBER_SECRET_KEY"))
	if len(key) != 32 {
		return "", errors.New("secret key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plainID), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptIDNumber decrypts a base64-encoded AES-GCM ciphertext back to plaintext
func DecryptIDNumber(encryptedID string) (string, error) {
	key := []byte(os.Getenv("IDNUMBER_SECRET_KEY"))
	if len(key) != 32 {
		return "", errors.New("secret key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedID)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
