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
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Convert smtpPort string to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	// Authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Sender and recipient
	from := username // Using the same username as sender
	to := []string{email}

	// Email content
	subject := "Your PIN Code"
	body := "Your PIN code is: " + pin

	// Constructing the email message
	message := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Sending the email
	err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, from, to, message)
	if err != nil {
		return err
	}
	log.Println("PIN code sent successfully to", email)
	return nil
}
