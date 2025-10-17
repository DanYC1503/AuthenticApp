package encryption

import (
	"net/http"
)

// validateSessionToken checks the session token by calling the auth-service
func ValidateSessionToken(r *http.Request) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9999/auth/validateToken", nil)

	// Pass the headers along (Authorization, cookies, etc)
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// You could also parse the JSON body if your validateToken returns {"tokenStatus":"Token valid"} or "invalid"
	return resp.StatusCode == http.StatusOK
}

func ValidateDeleteToken(r *http.Request) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9999/auth/validateDelToken", nil)

	// Pass the headers along (Authorization, cookies, etc)
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

	// Pass the headers along (Authorization, cookies, etc)
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
