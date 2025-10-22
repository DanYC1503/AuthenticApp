package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var csrfSecret = []byte(os.Getenv("CSRF_SECRET"))

// Handler to issue CSRF token cookie
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := GenerateCSRFToken()

	http.SetCookie(w, &http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    token,
		HttpOnly: false, // Angular must read it
		Secure:   true,  // only over HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"csrf":"ok"}`)
}

// Generates token: <random>.<hmac>
func GenerateCSRFToken() string {
	randomBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		panic("failed to generate CSRF token")
	}

	randomPart := base64.RawURLEncoding.EncodeToString(randomBytes)
	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(randomPart))
	signature := mac.Sum(nil)
	sigPart := base64.RawURLEncoding.EncodeToString(signature)

	return fmt.Sprintf("%s.%s", randomPart, sigPart)
}

// ValidateCSRFRequest validates CSRF token using cookie + header from *http.Request
func ValidateCSRFRequest(r *http.Request) bool {
	cookie, err := r.Cookie("XSRF-TOKEN")
	if err != nil {
		return false
	}

	headerToken := r.Header.Get("X-XSRF-TOKEN")
	if headerToken == "" {
		headerToken = r.FormValue("csrf_token")
	}

	// Both must match exactly
	if headerToken == "" || cookie.Value == "" || headerToken != cookie.Value {
		return false
	}

	// Now validate HMAC integrity
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return false
	}

	randomPart, sigPart := parts[0], parts[1]
	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(randomPart))
	expectedSig := mac.Sum(nil)

	sig, err := base64.RawURLEncoding.DecodeString(sigPart)
	if err != nil {
		return false
	}

	return hmac.Equal(sig, expectedSig)
}
