package controllers

import (
	"encoding/json"
	"log"
	"main/config"
	"main/models"
	auditprocessing "main/repository/audit-processing"
	"net/http"
	"time"
)

func AuditActionHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[Audit] Panic: %v", rec)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	const maxRetries = 3

	// Only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body FIRST
	var audit models.AuditLog
	if err := json.NewDecoder(r.Body).Decode(&audit); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// If not sent in body, check headers (fallback)
	if audit.Username == "" {
		audit.Username = r.Header.Get("X-Username")
	}
	if audit.Action == "" {
		audit.Action = r.Header.Get("X-Action")
	}
	if audit.IPAddress == "" {
		audit.IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if audit.UserAgent == "" {
		audit.UserAgent = r.Header.Get("User-Agent")
	}
	if audit.Timestamp.IsZero() {
		audit.Timestamp = time.Now().UTC()
	}

	// Validate required fields
	if audit.Username == "" || audit.Action == "" {
		http.Error(w, "Missing required fields: username and action", http.StatusBadRequest)
		return
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		db := config.ConnectDB()

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Retry %d/%d: Begin failed: %v\n", attempt, maxRetries, err)
			db.Close()
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		err = auditprocessing.AuditAction(tx, audit)
		if err != nil {
			// Rollback immediately on error
			tx.Rollback()
			db.Close()
			log.Printf("Retry %d/%d: AuditAction failed: %v\n", attempt, maxRetries, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Retry %d/%d: Commit failed: %v\n", attempt, maxRetries, err)
			db.Close()
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		db.Close()

		// Success - send response
		resp := map[string]string{"status": "logged"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("[Audit] CRITICAL: Failed to log audit after %d retries: %v", maxRetries, lastErr)
	http.Error(w, "Failed to log action after retries", http.StatusInternalServerError)
}
