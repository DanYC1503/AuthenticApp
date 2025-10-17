package main

import (
	"log"
	internal "main/internal" // internal package with routes & proxy
	"net/http"
	"os"
)

func main() {
	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	// Create main mux
	mux := http.NewServeMux()

	// Register routes (no CSRF)
	internal.RegisterRoutes(mux)

	log.Printf("API Gateway running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
