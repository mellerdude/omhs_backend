package utils

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ProjectManager struct{}

func NewProjectManager() *ProjectManager {
	return &ProjectManager{}
}

func (pm *ProjectManager) Execute(fn func() error, msg string) {
	if err := fn(); err != nil {
		pm.HandleError(err, msg)
	}
}

func (pm *ProjectManager) HandleError(err error, msg string) {
	if err == nil {
		return
	}
	fatal := false
	for _, kw := range []string{"fatal", "connect", "ping"} {
		if strings.Contains(strings.ToLower(msg), kw) {
			fatal = true
			break
		}
	}
	if fatal {
		logrus.Fatalf("%s: %v", msg, err)
	} else {
		logrus.Errorf("%s: %v", msg, err)
	}
}
