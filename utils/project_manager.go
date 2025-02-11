package utils

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once                   sync.Once
	projectManagerInstance *ProjectManager
)

// ProjectManager is a struct for managing project tasks and error handling
type ProjectManager struct{}

// NewProjectManager creates a new instance of ProjectManager
func NewProjectManager() *ProjectManager {
	once.Do(func() {
		projectManagerInstance = &ProjectManager{}
	})
	return projectManagerInstance
}

// Execute handles function execution and error logging
func (pm *ProjectManager) Execute(fn func() (interface{}, error), errorMessage string) {
	_, err := fn()
	if err != nil {
		pm.HandleError(err, errorMessage)
	}
}

// HandleError checks error and logs appropriately
func (pm *ProjectManager) HandleError(err error, errorMessage string) {
	if err != nil {
		fatalConditions := []string{"fatal", "connect", "ping"}
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
