package internal

import (
	"main/handlers"
	middleware "main/middleware/csrf"
	"net/http"
	"os"
)

// RegisterRoutes sets up routes and proxies for microservices.
func RegisterRoutes(mux *http.ServeMux) {
	authService := os.Getenv("AUTH_SERVICE_URL")
	if authService == "" {
		authService = "http://localhost:9999"
	}

	userService := os.Getenv("USER_SERVICE_URL")
	if userService == "" {
		userService = "http://localhost:8889"
	}
	auditService := os.Getenv("AUDIT_SERVICE_URL")
	if auditService == "" {
		auditService = "http://localhost:8890"
	}

	// Health check (unprotected)
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Gateway running"))
	})

	// Proxy routes
	mux.Handle("/auth/", handlers.ReverseProxy(authService, "/auth"))
	mux.Handle("/users/", handlers.ReverseProxy(userService, "/users"))
	mux.HandleFunc("/api/csrf-token", middleware.GetCSRFToken)
}
