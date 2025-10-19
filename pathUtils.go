package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetExecutablePath returns the directory where the executable is located
func GetExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return filepath.Dir(execPath), nil
}

// GetExecutableRelativePath returns a path relative to the executable directory
func GetExecutableRelativePath(relativePath string) (string, error) {
	execDir, err := GetExecutablePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(execDir, relativePath), nil
}
