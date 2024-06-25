package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	r_error "github.com/DanWlker/remind/error"
)

func FormatPathToRemoveHome(filePathWithHome string) (string, error) {
	home, errUserHomeDir := os.UserHomeDir()
	if errUserHomeDir != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", errUserHomeDir)
	}

	if !strings.HasPrefix(filePathWithHome, home) {
		return filePathWithHome, &r_error.FilePathNotStartsWithHome{
			HomeStr: home,
			FileStr: filePathWithHome,
		}
	}

	return strings.TrimPrefix(filePathWithHome, home), nil
}

func GetCurrentProgramExecutionDirectory() (string, error) {
	ex, errExecutable := os.Executable()
	if errExecutable != nil {
		return "", fmt.Errorf("os.Executable: %w", errExecutable)
	}

	return filepath.Dir(ex), nil
}

func GetHomeRemovedCurrentProgramExecutionDirectory() (string, error) {
	currProExDir, errGetCurrProExDir := GetCurrentProgramExecutionDirectory()
	if errGetCurrProExDir != nil {
		return "", fmt.Errorf("GetCurrentProgramExecutionDirectory: %w", errGetCurrProExDir)
	}

	path, errFormatPathToRemoveHome := FormatPathToRemoveHome(currProExDir)
	if errFormatPathToRemoveHome != nil {
		return "", fmt.Errorf("FormatPathToRemoveHome: %w", errFormatPathToRemoveHome)
	}

	return path, nil
}
