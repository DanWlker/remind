package constant

import (
	"os"
)

const DEFAULT_CONFIG_FILE_TYPE = "yaml"
const DEFAULT_CONFIG_FILE_NAME = ".remind"
const DEFAULT_CONFIG_FULL_FILE_NAME = string(os.PathSeparator) + DEFAULT_CONFIG_FILE_NAME + "." + DEFAULT_CONFIG_FILE_TYPE
const DEFAULT_CONFIG_PATH_AFTER_HOME = string(os.PathSeparator) + ".config"
const DEFAULT_CONFIG_FULL_PATH = DEFAULT_CONFIG_PATH_AFTER_HOME + DEFAULT_CONFIG_FULL_FILE_NAME

const DEFAULT_DATA_PATH_AFTER_HOME = string(os.PathSeparator) + "remind"
