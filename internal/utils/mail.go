package utils

import (
	"crypto/tls"
	"net/smtp"
	"os"

	"github.com/sirupsen/logrus"
)

func SendEmail(to, subject, message string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + message + "\r\n")

	c, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		logrus.Errorf("SMTP connect error: %v", err)
		return err
	}
	defer c.Close()

	if smtpPort == "587" {
		if err = c.StartTLS(&tls.Config{ServerName: smtpHost}); err != nil {
			return err
		}
	}

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

	logrus.Infof("Email sent successfully to %s", to)
	return nil
}
