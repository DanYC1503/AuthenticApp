package middleware

import (
	"net/http"
)

func LogoutCurrentUser(w http.ResponseWriter, r *http.Request) bool {
	// Remove the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Could add checks here if needed; for now always return 1
	return true
}
