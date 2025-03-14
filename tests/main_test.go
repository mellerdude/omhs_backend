package tests

import (
	"os"
	"testing"
)

// Global test managers for multiple suites
var authTestManager *TestManager
var requestsTestManager *TestManager
var mailTestManager *TestManager

func TestMain(m *testing.M) {
	// Initialize test managers for each suite
	authTestManager = GetTestManager("auth_test suite")
	requestsTestManager = GetTestManager("requests_test suite")
	mailTestManager = GetTestManager("mail_test suite")

	// Run all tests
	exitCode := m.Run()

	// Print the summaries for each suite after all tests are run
	authTestManager.PrintSummary()
	requestsTestManager.PrintSummary()
	mailTestManager.PrintSummary()

	PrintOverallSummary()

	// Exit with the appropriate code
	os.Exit(exitCode)
}
