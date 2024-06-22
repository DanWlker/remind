package helper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
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

func GetDataFile() string {
	folder, errGetDataFolder := GetDataFolder()
	if errGetDataFolder != nil {
		log.Println(errGetDataFolder)
	}
	dataFile := folder + string(os.PathSeparator) + "tempdata.yaml"

	_, errStat := os.Stat(dataFile)
	if errors.Is(errStat, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(dataFile), 0770)
		if err != nil {
			debug.PrintStack()
			log.Println(err)
		}
		if _, errFileCreate := os.Create(dataFile); errFileCreate != nil {
			debug.PrintStack()
			log.Println(err)
		}
	} else if errStat != nil {
		debug.PrintStack()
		log.Println(errStat)
	}

	return dataFile
}
