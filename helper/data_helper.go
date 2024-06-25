package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	r_error "github.com/DanWlker/remind/error"
	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

// This does not create the file if it doesn't exist
func GetTodoFromDataFile(dataFileFullPath string) ([]entity.TodoEntity, error) {
	file, errReadFile := os.ReadFile(dataFileFullPath)
	if errReadFile != nil {
		return []entity.TodoEntity{}, fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []entity.TodoEntity
	if errUnmarshal := yaml.Unmarshal(file, &items); errUnmarshal != nil {
		return []entity.TodoEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil
}

func WriteTodoToFile(fileFullPath string) error {
	return nil
}

func GetDataFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(constant.DATA_FOLDER_KEY))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + constant.DEFAULT_DATA_PATH_AFTER_HOME
	}

	if errMkDirAll := os.MkdirAll(dataFolder, 0770); errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return dataFolder, nil
}

// func GetDataFile(fileName string) (string, error) {
// 	dataFolder, errGetDataFolder := GetDataFolder()
// 	if errGetDataFolder != nil {
// 		return "", errGetDataFolder
// 	}
//
// 	fileFullPath := dataFolder + string(os.PathSeparator) + fileName
//
// 	_, errStat := os.Stat(fileFullPath)
// 	if errors.Is(errStat, os.ErrNotExist) {
// 		_, errCreate := os.Create(fileFullPath)
// 		if errCreate != nil {
// 			return "", fmt.Errorf("os.Create: %w", errCreate)
// 		}
// 	} else if errStat != nil {
// 		return "", errStat
// 	}
//
// 	return fileFullPath, nil
// }

func PrettyPrintDataFile(dataFileFullPath string, prefix string) error {
	todoList, errGetTodoFromDataFile := GetTodoFromDataFile(dataFileFullPath)
	if errGetTodoFromDataFile != nil {
		return fmt.Errorf("GetTodoFromDataFile: %w", errGetTodoFromDataFile)
	}

	for _, todo := range todoList {
		fmt.Println(prefix + todo.Text)
	}
	return nil
}

func GetRecordFile() (string, error) {
	dataFolder, errGetDataFolder := GetDataFolder()
	if errGetDataFolder != nil {
		return "", fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	defaultDataRecordFileFullPath := dataFolder + constant.DEFAULT_DATA_RECORD_FULL_FILE_NAME

	if _, errStat := os.Stat(defaultDataRecordFileFullPath); errors.Is(errStat, os.ErrNotExist) {
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

func FindProjectRecordFromFileWith(homeRemovedFolderPath string) (entity.ProjectRecordEntity, error) {
	allRecords, errGetRecordFileContents := GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return entity.ProjectRecordEntity{}, fmt.Errorf("GetRecordFileContents: %w", errGetRecordFileContents)
	}

	for _, record := range allRecords {
		if record.Path == homeRemovedFolderPath {
			return record, nil
		}
	}

	return entity.ProjectRecordEntity{}, fmt.Errorf("Record %v does not exst: %w", homeRemovedFolderPath, &r_error.RecordDoesNotExistError{})
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
