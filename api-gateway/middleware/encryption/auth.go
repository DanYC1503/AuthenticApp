package encryption

import (
	"fmt"
	"log"
	"net/http"
)

// validateSessionToken checks the session token by calling the auth-service endpoint
func ValidateSessionToken(r *http.Request) bool {
	fmt.Println("Trying to validate sessionToken for user")

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9999/auth/validateToken", nil)

	req.Header = r.Header

	// Send request to auth-service
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error contacting auth-service:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func ValidateDeleteToken(r *http.Request) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9999/auth/validateDelToken", nil)

	// Pass the headers (Authorization, cookies, etc)
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
func ValidateUpdateToken(r *http.Request) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9999/auth/validateUpToken", nil)

	// Pass the headers (Authorization, cookies, etc)
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
