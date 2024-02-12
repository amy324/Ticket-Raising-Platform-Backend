package main

import (
	// "fmt"
	// "log"

	"github.com/go-gomail/gomail"
)

// sendEmail sends an email using the provided parameters
func sendPinByEmail(recipient, subject, body string) error {
    // SMTP configuration
    smtpHost := "localhost" // MailHog SMTP host
    smtpPort := 1025        // MailHog SMTP port

    // Sender
    sender := "ticketplatform@email.com"

    // Compose the email message
    m := gomail.NewMessage()
    m.SetHeader("From", sender)
    m.SetHeader("To", recipient)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    // Send the email
    d := gomail.NewDialer(smtpHost, smtpPort, "", "")


    // Send the email
    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}