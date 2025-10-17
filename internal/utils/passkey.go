package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GeneratePasskey() (string, error) {
	bytes := make([]byte, 4) // short key
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
