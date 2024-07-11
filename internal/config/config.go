package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Keys for config
const (
	UserDefinedDataFolder = "UserDefinedDataFolder"
	ListAfterAddingTodo   = "ListAfterAddingTodo"
	ListAfterEditingTodo  = "ListAfterEditingTodo"
	ListAfterRemovingTodo = "ListAfterRemovingTodo"
)

// Data Defaults
const DefaultDataSubdir = "remind"

const (
	defaultDataRecordFileType     = "yaml"
	defaultDataRecordFileName     = ".rrecord"
	DefaultDataRecordFullFileName = defaultDataRecordFileName + "." + defaultDataRecordFileType
)

const (
	defaultDataFileFileType      = "yaml"
	DefaultDataFileFileExtension = ".rdata" + "." + defaultDataFileFileType
)

// Config Defaults
const (
	DefaultConfigFileType = "yaml"
	DefaultConfigFileName = ".remind"
)

func GetConfigFolder() (string, error) {
	// This has a different implementation from the data version as it is
	// impossible to get user defined config file location when it has not
	// been initialized at all
	configFolder, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("os.UserConfigDir: %w", err)
	}

	configFolder = filepath.Join(configFolder, "remind")

	if err := os.MkdirAll(configFolder, 0o770); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	return configFolder, nil
}
