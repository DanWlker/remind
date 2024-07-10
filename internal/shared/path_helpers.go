package shared

import (
	"fmt"
	"os"
	"strings"

	i_error "github.com/DanWlker/remind/internal/error"
)

func FormatRemoveHome(filePathWithHome string) (string, error) {
	home, errUserHomeDir := os.UserHomeDir()
	if errUserHomeDir != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", errUserHomeDir)
	}

	if !strings.HasPrefix(filePathWithHome, home) {
		return filePathWithHome, &i_error.NotUnderHomeError{
			Home: home,
			File: filePathWithHome,
		}
	}

	return strings.TrimPrefix(filePathWithHome, home), nil
}

func GetHomeRemovedWorkingDir() (string, error) {
	currProExDir, errGetwd := os.Getwd()
	if errGetwd != nil {
		return "", fmt.Errorf("Getwd: %w", errGetwd)
	}

	path, errFormatPathToRemoveHome := FormatRemoveHome(currProExDir)
	if errFormatPathToRemoveHome != nil {
		return "", fmt.Errorf("FormatPathToRemoveHome: %w", errFormatPathToRemoveHome)
	}

	return path, nil
}
