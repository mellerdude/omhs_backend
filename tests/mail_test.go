package tests

import (
	"omhs-backend/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSendEmail tests the sendEmail function
func TestSendEmail(t *testing.T) {
	// Call the sendEmail function
	err := utils.SendEmail("recipient@example.com", "Test Subject", "Test Message")
	assert.NoError(t, err)

	mailTestManager.RegisterTest(t, "TestSendEmail")
}
