package helper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	r_error "github.com/DanWlker/remind/error"
	"github.com/goccy/go-yaml"
)

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
		globalRecordEntity, errCreateNewRecord := CreateNewRecord("")
		if errCreateNewRecord != nil {
			return "", fmt.Errorf("CreateNewRecord: %w", errCreateNewRecord)
		}

		if err := SetRecordFileContents([]entity.ProjectRecordEntity{globalRecordEntity}); err != nil {
			return "", fmt.Errorf("SetRecordFileContents: %w", err)
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

	return entity.ProjectRecordEntity{}, &r_error.RecordDoesNotExistError{
		RecordIdentifier: homeRemovedFolderPath,
	}
}
func CreateNewRecord(pathIdentifier string) (entity.ProjectRecordEntity, error) {
	dataFolder, errGetDataFolder := GetDataFolder()
	if errGetDataFolder != nil {
		return entity.ProjectRecordEntity{}, fmt.Errorf("GetDataFolder: %w", errGetDataFolder)
	}

	newFile, errCreateTemp := os.CreateTemp(dataFolder, "*"+constant.DEFAULT_DATA_FILE_EXTENSION)
	if errCreateTemp != nil {
		return entity.ProjectRecordEntity{}, fmt.Errorf("os.CreateTemp: %w", errCreateTemp)
	}

	_, fileName := filepath.Split(newFile.Name())
	currentDirectoryRecord := entity.ProjectRecordEntity{
		DataFileName: fileName,
		Path:         pathIdentifier,
	}
	return currentDirectoryRecord, nil
}
