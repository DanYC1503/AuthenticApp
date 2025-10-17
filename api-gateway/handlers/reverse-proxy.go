package handlers

import (
	"log"
	"main/middleware/encryption"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ReverseProxy(target string, prefix string) http.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Host = targetURL.Host

		path := r.URL.Path

		// Protect /users/ and /payments/ routes
		if strings.HasPrefix(path, "/users/") || strings.HasPrefix(path, "/payments/") {
			if !encryption.ValidateSessionToken(r) {
				http.Error(w, "Unauthorized: invalid session token", http.StatusUnauthorized)
				return
			}

			// Extra validation for DELETE /users/delUser
			if path == "/users/delUser" && r.Method == http.MethodDelete {
				if !encryption.ValidateDeleteToken(r) {
					http.Error(w, "Unauthorized: invalid delete token", http.StatusUnauthorized)
					return
				}
			}

			// Extra validation for PUT /users/update
			if path == "/users/update" && r.Method == http.MethodPut {
				if !encryption.ValidateUpdateToken(r) {
					http.Error(w, "Unauthorized: invalid update token", http.StatusUnauthorized)
					return
				}
			}
		}

		log.Printf("Forwarding %s %s → %s%s", r.Method, path, targetURL, path)
		proxy.ServeHTTP(w, r)
	}
}
