package middleware

import (
	"fmt"
	"os"
)

// Email constructor with redirect if verified to frontend endpoint to pass through the api, verify and recieve the password token change
func SendPasswordRecoveryEmail(toEmail, token string) error {

	appURL := os.Getenv("FRONTEND_URL")
	if appURL == "" {
		return fmt.Errorf("FRONTEND_URL not set in environment")
	}

	// Build the reset link pointing to the Angular route
	resetLink := fmt.Sprintf("%s/passwordRecovery?token=%s", appURL, token)


	subject := "Password Reset Request"
	body := fmt.Sprintf(`Hello,

We received a request to reset your password :D.

Click the link below to reset your password:
%s

This link will expire in 15 minutes.

If you did not request this, you can safely ignore it.
`, resetLink)

	return SendEmail(toEmail, subject, body)
}
