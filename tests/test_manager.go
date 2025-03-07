// test_manager.go
package tests

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// TestManager handles a single test suite
type TestManager struct {
	suiteName   string
	passedCount int32
	failedCount int32
}

var (
	managers = make(map[string]*TestManager) // Store multiple test managers by suite name
)

// GetTestManager retrieves or creates a TestManager for the given suite
func GetTestManager(suiteName string) *TestManager {
	if _, exists := managers[suiteName]; !exists {
		managers[suiteName] = &TestManager{suiteName: suiteName}
	}
	return managers[suiteName]
}

// RegisterTest tracks the result of a test
func (tm *TestManager) RegisterTest(t *testing.T, testName string) {
	if t.Failed() {
		tm.failedCount++
		logrus.Infof("Test %s failed.", testName)
	} else {
		tm.passedCount++
		logrus.Infof("Test %s passed.", testName)
	}
}

// PrintSummary prints the test results
func (tm *TestManager) PrintSummary() {
	totalTests := tm.passedCount + tm.failedCount
	logrus.Infof("[%s] Tests Summary: Tests: %d Passed: %d Failed: %d",
		tm.suiteName, totalTests, tm.passedCount, tm.failedCount)
}
