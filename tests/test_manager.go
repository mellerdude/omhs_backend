package tests

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type TestManager struct {
	tests       []func(*testing.T)
	testNames   []string
	suiteName   string
	passedCount int
	failedCount int
}

func NewTestManager(suiteName string) *TestManager {
	return &TestManager{suiteName: suiteName}
}

func (tm *TestManager) AddTest(testFunc func(*testing.T)) {
	funcName := runtime.FuncForPC(reflect.ValueOf(testFunc).Pointer()).Name()
	// Clean up the function name to remove the package path
	funcNameParts := strings.Split(funcName, ".")
	cleanFuncName := funcNameParts[len(funcNameParts)-1]
	tm.tests = append(tm.tests, testFunc)
	tm.testNames = append(tm.testNames, cleanFuncName)
}

func (tm *TestManager) AddTests(testFuncs ...func(*testing.T)) {
	for _, testFunc := range testFuncs {
		tm.AddTest(testFunc)
	}
}

func (tm *TestManager) RunTests() {
	tm.logSuiteStart()
	for i, test := range tm.tests {
		logrus.Infof("Starting Test %d: %s", i+1, tm.testNames[i])
		t := &testing.T{}
		test(t)
		if !t.Failed() {
			logrus.Infof("Test %d: %s passed successfully.", i+1, tm.testNames[i])
			tm.passedCount++
		} else {
			logrus.Infof("Test %d: %s failed.", i+1, tm.testNames[i])
			tm.failedCount++
		}
		logrus.Infof("Completed Test %d: %s", i+1, tm.testNames[i])
	}
	tm.PrintSummary()
	tm.logSuiteEnd()
}

func (tm *TestManager) logSuiteStart() {
	logrus.Infof("Starting test suite: %s", tm.suiteName)
}

func (tm *TestManager) logSuiteEnd() {
	logrus.Infof("Completed test suite: %s", tm.suiteName)
}

func (tm *TestManager) PrintSummary() {
	logrus.Infof("Tests Summary: Tests: %d Passed: %d Failed: %d", len(tm.tests), tm.passedCount, tm.failedCount)
}
