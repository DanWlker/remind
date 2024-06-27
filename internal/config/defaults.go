package config

import (
	"os"
)

// CONFIG
const DEFAULT_CONFIG_PATH_AFTER_HOME = string(os.PathSeparator) + ".config"

const DEFAULT_CONFIG_FILE_TYPE = "yaml"
const DEFAULT_CONFIG_FILE_NAME = ".remind"
const DEFAULT_CONFIG_FULL_FILE_NAME = string(os.PathSeparator) + DEFAULT_CONFIG_FILE_NAME + "." + DEFAULT_CONFIG_FILE_TYPE

// DATA
const DEFAULT_DATA_PATH_AFTER_HOME = string(os.PathSeparator) + "remind"

const DEFAULT_DATA_RECORD_FILE_TYPE = "yaml"
const DEFAULT_DATA_RECORD_FILE_NAME = ".rrecord"
const DEFAULT_DATA_RECORD_FULL_FILE_NAME = string(os.PathSeparator) + DEFAULT_DATA_RECORD_FILE_NAME + "." + DEFAULT_DATA_RECORD_FILE_TYPE

const DEFAULT_DATA_FILE_FILE_TYPE = "yaml"
const DEFAULT_DATA_FILE_EXTENSION = ".rdata" + "." + DEFAULT_DATA_FILE_FILE_TYPE
