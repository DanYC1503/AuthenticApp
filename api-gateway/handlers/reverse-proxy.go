package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	middleware "main/middleware/csrf"
	"main/middleware/encryption"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// responseRecorder lets us capture the response status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if rr.body != nil {
		rr.body.Write(b)
	}
	return rr.ResponseWriter.Write(b)
}

func ReverseProxy(target string, prefix string) http.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		// Preserve headers
		for k, vv := range req.Header {
			req.Header[k] = vv
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !(r.Method == http.MethodGet || r.Method == http.MethodHead || path == "/api/csrf-token") {
			if !middleware.ValidateCSRFRequest(r) {
				http.Error(w, "Forbidden: invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		// Debug: log cookies and headers
		for _, c := range r.Cookies() {
			log.Printf("[Proxy] Cookie: %s=%s", c.Name, c.Value)
		}

		// Protect /users/ routes
		if strings.HasPrefix(path, "/users/") {
			if !encryption.ValidateSessionToken(r) {
				http.Error(w, "Unauthorized: invalid session token", http.StatusUnauthorized)
				return
			}

			if path == "/users/delUser" && r.Method == http.MethodDelete {
				if !encryption.ValidateDeleteToken(r) {
					http.Error(w, "Unauthorized: invalid delete token here", http.StatusUnauthorized)
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

		// --- Copy body for both proxy and audit logging ---
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[Audit] Failed to read body: %v", err)
			bodyBytes = []byte{}
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           nil, // Set to bytes.Buffer if you need response body
		}

		proxy.ServeHTTP(rec, r)

		// Only audit successful requests (2xx status codes)
		if rec.statusCode >= 200 && rec.statusCode < 300 {
			LogAuditAction(r, path, bodyBytes)
		}
	}
}

func LogAuditAction(r *http.Request, path string, bodyBytes []byte) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = "unknown"
	}
	ua := r.Header.Get("User-Agent")
	if ua == "" {
		ua = "unknown"
	}

	username := ""
	if len(bodyBytes) > 0 {
		var payload encryption.RequestPayload
		if err := json.Unmarshal(bodyBytes, &payload); err == nil {
			username = payload.Username
		}
	}

	// Non-blocking send to queue
	select {
	case auditQueue <- auditEvent{
		username: username,
		ip:       ip,
		ua:       ua,
		method:   r.Method,
		path:     path,
	}:
		// Successfully queued
	default:
		// Queue is full, log error but don't block
		log.Printf("[Audit] WARNING: Audit queue full, dropping event")
	}
}
