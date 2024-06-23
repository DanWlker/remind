package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/spf13/viper"
)

func GetDataFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(constant.DATA_FOLDER_KEY))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + constant.DEFAULT_DATA_PATH_AFTER_HOME
	}

	// TODO: We can optimize this to only be checked if hit exception outside
	_, errStat := os.Stat(dataFolder)

	// switch {
	// case errors.Is(errStat, os.ErrNotExist):
	// 	errMkDirAll := os.MkdirAll(dataFolder, 0770)
	// 	if errMkDirAll != nil {
	// 		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	// 	}
	// case errStat != nil:
	// 	return "", fmt.Errorf("os.Stat: %w", errStat)
	// }

	if errors.Is(errStat, os.ErrNotExist) {
		errMkDirAll := os.MkdirAll(dataFolder, 0770)
		if errMkDirAll != nil {
			return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
		}
	} else if errStat != nil {
		return "", fmt.Errorf("os.Stat: %w", errStat)
	}

	return dataFolder, nil
}
