package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"main/config"
	repository "main/repository/Audit-processing"
	"net/http"
	"time"
)

func GetUserAuditActions(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	var lastErr error

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db := config.ConnectDB()
		defer db.Close()
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Retry %d/%d: Begin failed: %v\n", attempt, maxRetries, err)
			db.Close()
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Call the function and get the audit logs
		audits, err := repository.UserAuditActions(tx, email)
		if err != nil {
			tx.Rollback()
			db.Close()
			log.Printf("Retry %d/%d: UserAuditActions failed: %v\n", attempt, maxRetries, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if err := tx.Commit(); err != nil {
			db.Close()
			log.Printf("Retry %d/%d: Commit failed: %v\n", attempt, maxRetries, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Success - send retrieved audit logs as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(audits); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	// If all retries failed
	http.Error(w, fmt.Sprintf("Failed to retrieve audits: %v", lastErr), http.StatusInternalServerError)
}
func RetriveUsers(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	var lastErr error

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db := config.ConnectDB()
		if db == nil {
			lastErr = fmt.Errorf("failed to connect to DB")
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		users, err := repository.GetUsers(db, email) 
		db.Close()
		if err != nil {
			log.Printf("Retry %d/%d: GetUsers failed: %v\n", attempt, maxRetries, err)
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Success - send retrieved users as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	// If all retries failed
	http.Error(w, fmt.Sprintf("Failed to retrieve users: %v", lastErr), http.StatusInternalServerError)
}
