package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken creates a random 32-byte hex token.
// Used for auth tokens or API keys.
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
