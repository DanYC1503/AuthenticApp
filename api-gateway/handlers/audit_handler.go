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

	action := detectAction(event.path)
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

// detectAction determines the action type based on the path
func detectAction(path string) string {
	// Auth actions
	if strings.Contains(path, "/auth/login") {
		return "login"
	}

	// Logout actions
	if strings.Contains(path, "/auth/logout") {
		return "logout-user"
	}
	if strings.Contains(path, "/auth/register") {
		return "user-create"
	}
	if strings.Contains(path, "/auth/deleteToken") {
		return "delete-token-granted"
	}

	// User management

	if strings.Contains(path, "/users/delUser") {
		return "user_delete"
	}
	if strings.Contains(path, "/users/update") {
		return "user_update"
	}
	if strings.Contains(path, "/users/disable/user") {
		return "user_disabled"
	}
	if strings.Contains(path, "/users/enable/user") {
		return "user_enabled"
	}
	if strings.Contains(path, "/users/") {
		return "user_access"
	}

	// Data operations
	if strings.Contains(path, "/data/") {
		return "data_access"
	}
	if strings.Contains(path, "/files/") {
		return "file_access"
	}

	// Admin operations
	if strings.Contains(path, "/admin/") {
		return "admin_access"
	}

	// Default to generic access
	return "access"
}
