// mailer.go

package main

import (
	"log"
	"net/smtp"
	"os"
	"strconv"
)

// sendPinByEmail sends a PIN code to the specified email address
func sendPinByEmail(email, pin string) error {
	// SMTP server configuration
	smtpHost := os.Getenv("SMTP_HOST")     // SMTP server host
	smtpPortStr := os.Getenv("SMTP_PORT")  // SMTP server port as string
	username := os.Getenv("SMTP_USERNAME") // SMTP username for authentication
	password := os.Getenv("SMTP_PASSWORD") // SMTP password for authentication

	// Convert smtpPort string to integer
	smtpPort, err := strconv.Atoi(smtpPortStr) // Convert port string to integer
	if err != nil {
		return err // Return error if conversion fails
	}

	// Authentication using PlainAuth method
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Sender and recipient email addresses
	from := username      // Using the same username as sender
	to := []string{email} // Recipient email address

	// Email content
	subject := "Your PIN Code"         // Email subject
	body := "Your PIN code is: " + pin // Email body containing the PIN code

	// Constructing the email message
	message := []byte("From: ticketplatform@email.com\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Sending the email using SendMail method
	err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, from, to, message)
	if err != nil {
		return err // Return error if sending email fails
	}

	// Log message if the PIN code is sent successfully
	log.Println("PIN code sent successfully to", email)
	return nil // Return nil (no error) if everything succeeds
}
