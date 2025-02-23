package tests

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTesting(t *testing.T) {
	logrus.Info("Starting TestSomething")

	assert.True(t, true, "True is true!")

	logrus.Info("TestTesting completed successfully")
}
