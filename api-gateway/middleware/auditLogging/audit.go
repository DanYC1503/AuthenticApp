package auditlogging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var auditClient = &http.Client{
	Timeout: 5 * time.Second, // Prevent hanging requests
}

func LogAction(username, ip, ua, method, path string) error {
	auditService := os.Getenv("AUDIT_SERVICE_URL")
	if auditService == "" {
		auditService = "http://localhost:8890/audit/log"
	}

	action := detectAction(path)
	payload := map[string]string{
		"username":   username,
		"action":     action,
		"ip_address": ip,
		"user_agent": ua,
		"metadata":   "{}",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Audit] Error marshaling audit payload: %v", err)
		return err
	}

	// Retry logic with exponential backoff
	var lastErr error
	for i := 0; i < 3; i++ { // Retry up to 3 times
		resp, err := auditClient.Post(auditService, "application/json", bytes.NewBuffer(body))
		if err != nil {
			lastErr = err
			log.Printf("[Audit] Attempt %d: Error sending audit log: %v", i+1, err)
			time.Sleep(time.Duration(i*i) * 100 * time.Millisecond) // Exponential backoff
			continue
		}

		// Handle non-200 responses
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("audit service returned status: %d", resp.StatusCode)
			log.Printf("[Audit] Attempt %d: Unexpected status: %d", i+1, resp.StatusCode)
			resp.Body.Close()
			time.Sleep(time.Duration(i*i) * 100 * time.Millisecond)
			continue
		}

		resp.Body.Close()
		log.Printf("[Audit] Logged %s (%s) → %d", action, username, resp.StatusCode)
		return nil // Success
	}

	log.Printf("[Audit] Failed after 3 attempts: %v", lastErr)
	return lastErr
}

// You can map your routes to friendly action names here
func detectAction(path string) string {
	switch path {
	// ==== AUTH ROUTES ====
	case "/auth/register":
		return "user_register"
	case "/auth/login":
		return "user_login"
	case "/auth/logout":
		return "user_logout"

	case "/auth/google/login":
		return "google_oauth_login"
	case "/auth/google/callback":
		return "google_oauth_callback"

	case "/auth/validateToken":
		return "validate_token"
	case "/auth/validateUpToken":
		return "validate_update_token"
	case "/auth/validateDelToken":
		return "validate_delete_token"

	case "/auth/deleteToken":
		return "get_delete_token"
	case "/auth/updateUserToken":
		return "get_update_token"

	// ==== USER ROUTES ====
	case "/users/update":
		return "user_update"
	case "/users/info":
		return "get_user_info"
	case "/users/delete":
		return "user_delete"

	default:
		return "unknown_action"
	}
}
