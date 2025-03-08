package controllers

import (
	"crypto/tls"
	"net/smtp"
	"os"

	"github.com/sirupsen/logrus"
)

// sendEmail sends an email with the specified subject and message
// Parameters:
// - to: string - the recipient's email address
// - subject: string - the subject of the email
// - message: string - the body of the email
// Returns:
// - error: any error encountered during sending the email
func SendEmail(to, subject, message string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	// Connect to the server
	c, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		logrus.Errorf("Failed to connect to SMTP server: %v", err)
		return err
	}
	defer c.Close()

	// Start TLS if using port 587
	if smtpPort == "587" {
		if err = c.StartTLS(&tls.Config{ServerName: smtpHost}); err != nil {
			return err
		}
	}

	// Perform SMTP authentication
	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(msg); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	logrus.Infof("Email sent successfully to: %s", to)
	return nil
}
