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

	return true
}
