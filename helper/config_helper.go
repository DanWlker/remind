package helper

import (
	"fmt"
	"os"

	"github.com/DanWlker/remind/constant"
)

func GetConfigFolder() (string, error) {
	// This has a different implementation from the data version as it is
	// impossible to get user defined config file location when it has not
	// been initialized at all
	home, errHomeDir := os.UserHomeDir()
	if errHomeDir != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
	}
	configFolder := home + constant.DEFAULT_CONFIG_PATH_AFTER_HOME

	errMkDirAll := os.MkdirAll(configFolder, os.ModeDir)
	if errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return configFolder, nil
}
