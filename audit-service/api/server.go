package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// EnableServer starts the Auth service
func EnableServer() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8890"
	}

	// Create mux and register routes
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	// Configure and start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Audit service running on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
