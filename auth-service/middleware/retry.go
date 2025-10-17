package middleware

import (
	"database/sql"
	"log"
	"main/middleware/encryption"
	"net/http"
	"time"
)

func TransactionRetry(db *sql.DB, maxRetries int, w http.ResponseWriter, op func(tx *sql.Tx) error) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx, err := db.Begin()
		if err != nil {
			log.Println("Begin failed:", err)
			continue
		}

		err = op(tx) // Execute the operation/function passed as parameter

		if err != nil {
			tx.Rollback()
			if encryption.IsRetryable(err) {
				log.Printf("Retry %d/%d: %v\n", attempt, maxRetries, err)
				time.Sleep(time.Duration(attempt) * time.Second) // exponential backoff
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			if encryption.IsRetryable(err) {
				log.Printf("Retry %d/%d: Commit failed: %v\n", attempt, maxRetries, err)
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Transaction succeeded on attempt", attempt)
		w.WriteHeader(http.StatusCreated)
		break
	}
}
