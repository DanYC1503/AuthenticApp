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

	// Preserve original Host and headers for OAuth session
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		// Do NOT override req.Host, needed for goth/gothic session
		// req.Host = targetURL.Host

		// Preserve all headers, including cookies
		for k, vv := range req.Header {
			req.Header[k] = vv
		}

	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Debug: print cookies received by the proxy
		for _, c := range r.Cookies() {
			log.Printf("Cookie received: %s=%s", c.Name, c.Value)
		}

		// Protect /users/ routes
		if strings.HasPrefix(path, "/users/") {
			if !encryption.ValidateSessionToken(r) {
				http.Error(w, "Unauthorized: invalid session token users path", http.StatusUnauthorized)
				return
			}

			if path == "/users/delUser" && r.Method == http.MethodDelete {
				if !encryption.ValidateDeleteToken(r) {
					http.Error(w, "Unauthorized: invalid delete token", http.StatusUnauthorized)
					return
				}
			}

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
