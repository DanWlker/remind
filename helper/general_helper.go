package helper

import (
	"fmt"
	"os"
	"strings"

	r_error "github.com/DanWlker/remind/error"
	"github.com/spf13/cobra"
)

func GetHomeRemovedPath(fileFullPath string) (string, error) {
	home, errUserHomeDir := os.UserHomeDir()
	if errUserHomeDir != nil {
		cobra.CheckErr(fmt.Errorf("os.UserHomeDir: %w", errUserHomeDir))
	}

	if !strings.HasPrefix(fileFullPath, home) {
		return fileFullPath, &r_error.FilePathNotStartsWithHome{HomeStr: home}
	}

	return strings.TrimPrefix(fileFullPath, home), nil
}
