package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	"github.com/goccy/go-yaml"
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

func GetRecordFile() (string, error) {
	dataFolder, errGetDataFolder := GetDataFolder()
	if errGetDataFolder != nil {
		return "", fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	defaultDataRecordFileFullPath := dataFolder + constant.DEFAULT_DATA_RECORD_FULL_FILE_NAME

	_, errStat := os.Stat(defaultDataRecordFileFullPath)
	if errors.Is(errStat, os.ErrNotExist) {
		_, errCreate := os.Create(defaultDataRecordFileFullPath)
		if errCreate != nil {
			return "", fmt.Errorf("os.Create: %w", errCreate)
		}
	} else if errStat != nil {
		return "", fmt.Errorf("os.Stat: %w", errStat)
	}

	return defaultDataRecordFileFullPath, nil
}

func GetRecordFileContents() ([]entity.ProjectRecordEntity, error) {
	recordFileString, errGetRecordFile := GetRecordFile()
	if errGetRecordFile != nil {
		return []entity.ProjectRecordEntity{}, fmt.Errorf("GetRecordFile: %w", errGetRecordFile)
	}

	recordFile, errReadFile := os.ReadFile(recordFileString)
	if errReadFile != nil {
		return []entity.ProjectRecordEntity{}, fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []entity.ProjectRecordEntity
	if errUnmarshal := yaml.Unmarshal(recordFile, &items); errUnmarshal != nil {
		return []entity.ProjectRecordEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil

}

func SetRecordFileContents(items []entity.ProjectRecordEntity) error {
	recordFileString, errGetRecordFile := GetRecordFile()
	if errGetRecordFile != nil {
		return fmt.Errorf("GetRecordFile: %w", errGetRecordFile)
	}

	yamlContent, errMarshal := yaml.Marshal(items)
	if errMarshal != nil {
		return fmt.Errorf("yaml.Marshal: %w", errMarshal)
	}

	errWriteFile := os.WriteFile(recordFileString, yamlContent, 0644)
	if errWriteFile != nil {
		return fmt.Errorf("os.WriteFile: %w", errWriteFile)
	}

	return nil
}
