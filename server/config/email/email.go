package email

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendEmail sends an email with the provided subject and body
func SendEmail(to string, subject string, body string) error {
	// Load SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")     // e.g., "smtp.gmail.com"
	smtpPortStr := os.Getenv("SMTP_PORT")  // e.g., "587"
	smtpUsername := os.Getenv("SMTP_USER") // Your email username (e.g., "youremail@gmail.com")
	smtpPassword := os.Getenv("SMTP_PASS") // Your email password (or app password)
	from := os.Getenv("SMTP_FROM_EMAIL")   // From email address (e.g., "youremail@gmail.com")

	// Convert smtpPortStr to an integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Invalid SMTP_PORT: %v", err)
		return err
	}

	// Set up the email message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	// Set up the SMTP dialer (using the configuration values)
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	fmt.Println("Email sent successfully! ✨")
	return nil
}
