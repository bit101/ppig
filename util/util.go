package util

import (
	"errors"
	"os"
)

func DoesPathExist(projectPath string) bool {
	_, err := os.Stat(projectPath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
