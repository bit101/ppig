package util

import (
	"errors"
	"fmt"
	"os"
)

func DoesPathExist(projectPath string) bool {
	_, err := os.Stat(projectPath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func PrintRed(msg string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", msg)
}
