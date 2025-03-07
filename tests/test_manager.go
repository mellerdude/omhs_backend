package tests

import (
	"sync/atomic"
	"testing"

	"github.com/sirupsen/logrus"
)

type TestManager struct {
	suiteName   string
	passedCount int32
	failedCount int32
}

func NewTestManager(suiteName string) *TestManager {
	return &TestManager{suiteName: suiteName}
}

func (tm *TestManager) RegisterTest(t *testing.T, testName string) {
	logrus.Infof("Running test %s", testName)
	if t.Failed() {
		atomic.AddInt32(&tm.failedCount, 1)
		logrus.Infof("Test %s failed.", testName)
	} else {
		atomic.AddInt32(&tm.passedCount, 1)
		logrus.Infof("Test %s passed.", testName)
	}
}

func (tm *TestManager) PrintSummary() {
	totalTests := tm.passedCount + tm.failedCount
	logrus.Infof("Tests Summary: Tests: %d Passed: %d Failed: %d", totalTests, tm.passedCount, tm.failedCount)
}
