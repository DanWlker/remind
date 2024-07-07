package config

import (
	"fmt"
	"os"
)

// Keys for config
const USER_DEFINED_DATA_FOLDER = "USER_DEFINED_DATA_FOLDER"
const LIST_AFTER_ADDING_TODO = "LIST_AFTER_ADDING_TODO"
const LIST_AFTER_EDITING_TODO = "LIST_AFTER_EDITING_TODO"
const LIST_AFTER_REMOVING_TODO = "LIST_AFTER_REMOVING_TODO"

// Data Defaults
const DEFAULT_DATA_PATH_AFTER_HOME = string(os.PathSeparator) + "remind"

const DEFAULT_DATA_RECORD_FILE_TYPE = "yaml"
const DEFAULT_DATA_RECORD_FILE_NAME = ".rrecord"
const DEFAULT_DATA_RECORD_FULL_FILE_NAME = string(os.PathSeparator) + DEFAULT_DATA_RECORD_FILE_NAME + "." + DEFAULT_DATA_RECORD_FILE_TYPE

const DEFAULT_DATA_FILE_FILE_TYPE = "yaml"
const DEFAULT_DATA_FILE_EXTENSION = ".rdata" + "." + DEFAULT_DATA_FILE_FILE_TYPE

// Config Defaults
const DEFAULT_CONFIG_PATH_AFTER_HOME = string(os.PathSeparator) + ".config"

const DEFAULT_CONFIG_FILE_TYPE = "yaml"
const DEFAULT_CONFIG_FILE_NAME = ".remind"
const DEFAULT_CONFIG_FULL_FILE_NAME = string(os.PathSeparator) + DEFAULT_CONFIG_FILE_NAME + "." + DEFAULT_CONFIG_FILE_TYPE

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
