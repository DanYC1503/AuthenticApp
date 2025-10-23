package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Create a reusable HTTP client with timeouts
var auditClient = &http.Client{
	Timeout: 5 * time.Second,
}

var auditQueue = make(chan auditEvent, 1000) // Buffer 1000 events

type auditEvent struct {
	username string
	ip       string
	ua       string
	method   string
	path     string
	statusCode int
}

func init() {
	// Start background worker to process audit events
	go processAuditEvents()
}

func processAuditEvents() {
	for event := range auditQueue {
		// Implement retry logic here
		for i := 0; i < 3; i++ {
			err := sendAuditEvent(event)
			if err == nil {
				break // Success
			}
			if i == 2 { // Last attempt failed
				log.Printf("[Audit] CRITICAL: Failed to send audit event after 3 attempts: %v", err)
			}
			time.Sleep(time.Duration(i*i) * 100 * time.Millisecond)
		}
	}
}

func sendAuditEvent(event auditEvent) error {
	auditService := os.Getenv("AUDIT_SERVICE_URL")
	if auditService == "" {
		auditService = "http://localhost:8890/audit/log"
	}

	action := detectAction(event.path, event.statusCode)
	payload := map[string]string{
		"username":   event.username,
		"action":     action,
		"ip_address": event.ip,
		"user_agent": event.ua,
		"metadata":   "{}",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := auditClient.Post(auditService, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	log.Printf("[Audit] Logged %s (%s) → %d", action, event.username, resp.StatusCode)
	return nil
}
func detectAction(path string, statusCode int) string {
	path = strings.ToLower(path)

	var baseAction string

	// ==== AUTH ROUTES ====
	switch {
	case strings.Contains(path, "/auth/login"):
		baseAction = "user_login"
	case strings.Contains(path, "/auth/logout"):
		baseAction = "user_logout"
	case strings.Contains(path, "/auth/register"):
		baseAction = "user_register"
	case strings.Contains(path, "/auth/google/login"):
		baseAction = "google_oauth_login"
	case strings.Contains(path, "/auth/google/callback"):
		baseAction = "google_oauth_callback"
	case strings.Contains(path, "/auth/validateuptoken"):
		baseAction = "validate_update_token"
	case strings.Contains(path, "/auth/validatedeltoken"):
		baseAction = "validate_delete_token"
	case strings.Contains(path, "/auth/validatetoken"):
		baseAction = "validate_token"
	case strings.Contains(path, "/auth/deletetoken"):
		baseAction = "delete_token_granted"
	case strings.Contains(path, "/auth/updatetoken"):
		baseAction = "get_update_token"
	case strings.Contains(path, "/auth/passwordtoken"):
		baseAction = "recover_password_email_token_granted"
	case strings.Contains(path, "/auth/verifypasstoken"):
		baseAction = "verify_grant_password_recovery_token"

	// ==== USER ROUTES ====
	case strings.Contains(path, "/users/update"):
		baseAction = "user_update"
	case strings.Contains(path, "/users/info"):
		baseAction = "get_user_info"
	case strings.Contains(path, "/users/delete"), strings.Contains(path, "/users/deluser"):
		baseAction = "user_delete"
	case strings.Contains(path, "/users/disable/user"):
		baseAction = "user_disabled"
	case strings.Contains(path, "/users/enable/user"):
		baseAction = "user_enabled"
	case strings.Contains(path, "/users/"):
		baseAction = "user_access"
	default:
		baseAction = "unknown_action"
	}

	// ==== Add success/failure suffix ====
	if statusCode >= 200 && statusCode < 300 {
		return baseAction + "_success"
	}
	return baseAction + "_fail"
}
