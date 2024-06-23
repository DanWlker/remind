package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	r_error "github.com/DanWlker/remind/error"
)

func GetHomeRemovedPath(fileFullPath string) (string, error) {
	home, errUserHomeDir := os.UserHomeDir()
	if errUserHomeDir != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", errUserHomeDir)
	}

	if !strings.HasPrefix(fileFullPath, home) {
		return fileFullPath, &r_error.FilePathNotStartsWithHome{HomeStr: home}
	}

	return strings.TrimPrefix(fileFullPath, home), nil
}

func GetCurrentProgramExecutionDirectory() (string, error) {
	ex, errExecutable := os.Executable()
	if errExecutable != nil {
		return "", fmt.Errorf("os.Executable: %w", errExecutable)
	}
	return filepath.Dir(ex), nil
}
