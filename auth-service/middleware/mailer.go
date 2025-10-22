package middleware

import (
	"fmt"
	"net/smtp"
	"os"
)
//Mailer with teh env variables to send email with google smtp credentials
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if from == "" || password == "" {
		return fmt.Errorf("SMTP_EMAIL or SMTP_PASSWORD not set in environment")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	auth := smtp.PlainAuth("", from, password, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
