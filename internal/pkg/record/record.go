package record

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	i_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/pkg/data"

	"github.com/DanWlker/remind/internal/config"
	"github.com/goccy/go-yaml"
)

func GetRecordFile() (string, error) {
	dataFolder, errGetDataFolder := data.GetDataFolder()
	if errGetDataFolder != nil {
		return "", fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	defaultDataRecordFileFullPath := dataFolder + config.DEFAULT_DATA_RECORD_FULL_FILE_NAME

	if _, errStat := os.Stat(defaultDataRecordFileFullPath); errors.Is(errStat, os.ErrNotExist) {
		_, errCreate := os.Create(defaultDataRecordFileFullPath)
		if errCreate != nil {
			return "", fmt.Errorf("os.Create: %w", errCreate)
		}
		globalRecordEntity, errCreateNewRecord := CreateNewRecord("")
		if errCreateNewRecord != nil {
			return "", fmt.Errorf("CreateNewRecord: %w", errCreateNewRecord)
		}

		if err := SetRecordFileContents([]RecordEntity{globalRecordEntity}); err != nil {
			return "", fmt.Errorf("SetRecordFileContents: %w", err)
		}
	} else if errStat != nil {
		return "", fmt.Errorf("os.Stat: %w", errStat)
	}

	return defaultDataRecordFileFullPath, nil
}

func GetRecordFileContents() ([]RecordEntity, error) {
	recordFileString, errGetRecordFile := GetRecordFile()
	if errGetRecordFile != nil {
		return []RecordEntity{}, fmt.Errorf("GetRecordFile: %w", errGetRecordFile)
	}

	recordFile, errReadFile := os.ReadFile(recordFileString)
	if errReadFile != nil {
		return []RecordEntity{}, fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []RecordEntity
	if errUnmarshal := yaml.Unmarshal(recordFile, &items); errUnmarshal != nil {
		return []RecordEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil

}

func SetRecordFileContents(items []RecordEntity) error {
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

func GetProjectRecordFromFileWith(homeRemovedFolderPath string) (RecordEntity, error) {
	allRecords, errGetRecordFileContents := GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return RecordEntity{}, fmt.Errorf("GetRecordFileContents: %w", errGetRecordFileContents)
	}

	for _, record := range allRecords {
		if record.Path == homeRemovedFolderPath {
			return record, nil
		}
	}

	return RecordEntity{}, &i_error.RecordDoesNotExistError{
		RecordIdentifier: homeRemovedFolderPath,
	}
}

func CreateNewRecord(pathIdentifier string) (RecordEntity, error) {
	dataFolder, errGetDataFolder := data.GetDataFolder()
	if errGetDataFolder != nil {
		return RecordEntity{}, fmt.Errorf("GetDataFolder: %w", errGetDataFolder)
	}

	newFile, errCreateTemp := os.CreateTemp(dataFolder, "*"+config.DEFAULT_DATA_FILE_EXTENSION)
	if errCreateTemp != nil {
		return RecordEntity{}, fmt.Errorf("os.CreateTemp: %w", errCreateTemp)
	}

	_, fileName := filepath.Split(newFile.Name())
	currentDirectoryRecord := RecordEntity{
		DataFileName: fileName,
		Path:         pathIdentifier,
	}
	return currentDirectoryRecord, nil
}
