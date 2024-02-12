package main

import (
	"github.com/go-gomail/gomail"
)

func sendPinByEmail(recipient, subject, body string) error {
	// SMTP configuration
	smtpHost := "localhost" // MailHog SMTP host
	smtpPort := 8026        // MailHog SMTP port

	// Sender
	sender := "ticketplatform@email.com"

	// Compose the email message
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email using SMTP
	d := gomail.NewDialer(smtpHost, smtpPort, "", "")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
