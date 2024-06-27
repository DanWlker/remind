package config

import (
	"fmt"
	"os"
)

func GetConfigFolder() (string, error) {
	// This has a different implementation from the data version as it is
	// impossible to get user defined config file location when it has not
	// been initialized at all
	home, errHomeDir := os.UserHomeDir()
	if errHomeDir != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
	}
	configFolder := home + DEFAULT_CONFIG_PATH_AFTER_HOME

	if errMkDirAll := os.MkdirAll(configFolder, 0770); errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return configFolder, nil
}
