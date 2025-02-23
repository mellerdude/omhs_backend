package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GeneratePasskey() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Errorf("Failed to generate passkey: %v", err)
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}
