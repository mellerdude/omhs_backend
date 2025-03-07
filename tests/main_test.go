package tests

import (
	"os"
	"testing"
)

// Global test managers for multiple suites
var authTestManager *TestManager

var requestsTestManager *TestManager

func TestMain(m *testing.M) {
	// Initialize test managers for each suite
	authTestManager = GetTestManager("auth_test suite")
	requestsTestManager = GetTestManager("requests_test suite")

	// Run all tests
	exitCode := m.Run()

	// Print the summaries for each suite after all tests are run
	authTestManager.PrintSummary()
	requestsTestManager.PrintSummary()

	// Exit with the appropriate code
	os.Exit(exitCode)
}
