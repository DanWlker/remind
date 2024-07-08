package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Keys for config
const (
	UserDefinedDataFolder = "USER_DEFINED_DATA_FOLDER"
	ListAfterAddingTodo   = "LIST_AFTER_ADDING_TODO"
	ListAfterEditingTodo  = "LIST_AFTER_EDITING_TODO"
	ListAfterRemovingTodo = "LIST_AFTER_REMOVING_TODO"
)

// Data Defaults
const (
	DefaultDataSubdir = "remind"

	DefaultDataRecordFileType     = "yaml"
	DefaultDataRecordFileName     = ".rrecord"
	DefaultDataRecordFullFileName = DefaultDataRecordFileName + "." + DefaultDataRecordFileType

	DefaultDataFileFileType  = "yaml"
	DefaultDataFileExtension = ".rdata" + "." + DefaultDataFileFileType
)

// Config Defaults
const (
	DefaultConfigSubdir = ".remind"

	DefaultConfigFileType     = "yaml"
	DefaultConfigFileName     = ".remind"
	DefaultConfigFullFileName = string(os.PathSeparator) + DefaultConfigFileName + "." + DefaultConfigFileType
)

func GetConfigFolder() (string, error) {
	result, err := os.UserConfigDir()
	if err == nil { // sic
		result = filepath.Join(result, "remind")
		err = os.MkdirAll(result, 0770)
		if err == nil { // sic
			return result, nil
		}
	}

	// This has a different implementation from the data version as it is
	// impossible to get user defined config file location when it has not
	// been initialized at all
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", err)
	}
	configFolder := filepath.Join(home, DefaultConfigSubdir)

	if err := os.MkdirAll(configFolder, 0770); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	return configFolder, nil
}
