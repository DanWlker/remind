package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	"github.com/spf13/viper"
)

func ReadFromFile(fileFullName string) ([]entity.TodoEntity, error) {
	return []entity.TodoEntity{}, nil
}

func WriteToFile() {}

func GetDataFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(constant.DATA_FOLDER_KEY))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + constant.DEFAULT_DATA_PATH_AFTER_HOME
	}

	errMkDirAll := os.MkdirAll(dataFolder, 0770)
	if errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return dataFolder, nil
}

func GetDataFile(fileName string) (string, error) {
	dataFolder, errGetDataFolder := GetDataFolder()
	if errGetDataFolder != nil {
		return "", errGetDataFolder
	}

	fileFullPath := dataFolder + string(os.PathSeparator) + fileName

	_, errStat := os.Stat(fileFullPath)
	if errors.Is(errStat, os.ErrNotExist) {
		_, errCreate := os.Create(fileFullPath)
		if errCreate != nil {
			return "", fmt.Errorf("os.Create: %w", errCreate)
		}
	} else if errStat != nil {
		return "", errStat
	}

	return fileFullPath, nil
}
