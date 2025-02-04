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
func sendEmail(to, subject, message string) error {
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

	// Connect to the server, authenticate, set the sender and recipient, and send the email
	c, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		logrus.Errorf("Failed to connect to SMTP server: %v", err)
		return err
	}
	defer c.Close()

	// Start TLS if using port 587
	if smtpPort == "587" {
		if err = c.StartTLS(&tls.Config{ServerName: smtpHost}); err != nil {
			logrus.Errorf("Failed to start TLS: %v", err)
			return err
		}
	}

	// Perform SMTP authentication and log any issues
	if err = c.Auth(auth); err != nil {
		logrus.Errorf("SMTP authentication failed: %v", err)
		logrus.Errorf("Check if EMAIL_USER and EMAIL_PASS are correct. Current user: %s", from)
		return err
	}

	if err = c.Mail(from); err != nil {
		logrus.Errorf("Failed to set sender: %v", err)
		return err
	}

	if err = c.Rcpt(to); err != nil {
		logrus.Errorf("Failed to set recipient: %v", err)
		return err
	}

	w, err := c.Data()
	if err != nil {
		logrus.Errorf("Failed to send data: %v", err)
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		logrus.Errorf("Failed to write message: %v", err)
		return err
	}

	err = w.Close()
	if err != nil {
		logrus.Errorf("Failed to close writer: %v", err)
		return err
	}

	logrus.Infof("Email sent successfully to: %s", to)
	return nil
}
