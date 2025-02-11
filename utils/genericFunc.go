package utils

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// HandleError is a generic function to handle errors and log them
func HandleError(err error, errorMessage string) {
	if err != nil {
		// Define conditions for fatal errors
		fatalConditions := []string{"fatal", "connect", "ping"}

		// Check if the error message contains any fatal conditions
		isFatal := false
		for _, condition := range fatalConditions {
			if strings.Contains(strings.ToLower(errorMessage), condition) {
				isFatal = true
				break
			}
		}

		if isFatal {
			logrus.Fatalf("%s: %v", errorMessage, err)
		} else {
			logrus.Errorf("%s: %v", errorMessage, err)
		}
	}
}
