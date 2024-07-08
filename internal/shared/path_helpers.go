package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	i_error "github.com/DanWlker/remind/internal/error"
)

func FormatRemoveHome(filePathWithHome string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", err)
	}

	rel, err := filepath.Rel(home, filePathWithHome)
	if err != nil {
		return "", fmt.Errorf("computing relative path from %s to %s: %w", home, filePathWithHome, err)
	}

	if strings.HasPrefix(rel, "../") {
		return "", i_error.NotUnderHomeError{Home: home, File: filePathWithHome}
	}

	return rel, nil
}

func GetHomeRemovedWorkingDir() (string, error) {
	currProExDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("os.Getwd: %w", err)
	}

	path, err := FormatRemoveHome(currProExDir)
	if err != nil {
		return "", fmt.Errorf("FormatPathToRemoveHome: %w", err)
	}

	return path, nil
}
