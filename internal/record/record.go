package record

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DanWlker/remind/internal/config"
	"github.com/DanWlker/remind/internal/data"
	i_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/shared"
)

func GetFile() (string, error) {
	dataFolder, err := data.GetFolder()
	if err != nil {
		return "", fmt.Errorf("data.GetFolder: %w", err)
	}

	fullPath := filepath.Join(dataFolder, config.DefaultDataRecordFullFileName)

	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(fullPath)
		if err != nil {
			return "", fmt.Errorf("os.Create: %w", err)
		}
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return "", fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}

		globalRecordEntity, err := CreateNewRecord(homeRemoved)
		if err != nil {
			return "", fmt.Errorf("CreateNewRecord: %w", err)
		}

		if err := SetFileContents([]RecordEntity{globalRecordEntity}); err != nil {
			return "", fmt.Errorf("SetRecordFileContents: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("os.Stat: %w", err)
	}

	return fullPath, nil
}

func GetFileContents() (items []RecordEntity, err error) {
	recordFile, err := GetFile()
	if err != nil {
		return nil, fmt.Errorf("GetFile: %w", err)
	}

	f, err := os.Open(recordFile)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	items, err = shared.FGetStructFromYaml[RecordEntity](f)
	if err != nil {
		return nil, fmt.Errorf("FGetStructFromYaml: %w", err)
	}

	return items, nil
}

func SetFileContents(items []RecordEntity) (err error) {
	recordFile, err := GetFile()
	if err != nil {
		return fmt.Errorf("GetFile: %w", err)
	}

	f, err := os.OpenFile(recordFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	err = shared.FWriteStructToYaml[RecordEntity](f, items)
	if err != nil {
		return fmt.Errorf("shared.FWriteStructToYaml: %w", err)
	}

	return nil
}

func FGetRecordIdentityWithIdentifier(items []RecordEntity, identifier string) (RecordEntity, error) {
	for _, record := range items {
		if record.Path == identifier {
			return record, nil
		}
	}

	return RecordEntity{}, i_error.RecordDoesNotExistError{
		ID: identifier,
	}
}

func GetRecordEntityWithIdentifier(homeRemovedPath string) (RecordEntity, error) {
	allRecords, err := GetFileContents()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("GetFileContents: %w", err)
	}

	res, err := FGetRecordIdentityWithIdentifier(allRecords, homeRemovedPath)
	if err != nil {
		return RecordEntity{}, fmt.Errorf("FGetRecordIdentityWithIdentifier: %w", err)
	}

	return res, nil
}

func CreateNewRecord(pathIdentifier string) (RecordEntity, error) {
	dataFolder, err := data.GetFolder()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("data.GetFolder: %w", err)
	}

	// Note: This is persistent data, CreateTemp is used to help create
	// files that have names that don't collide with other files in the
	// folder. The first argument is dataFolder instead of the usual ""
	newFile, err := os.CreateTemp(dataFolder, "*"+config.DefaultDataFileFileExtension)
	if err != nil {
		return RecordEntity{}, fmt.Errorf("os.CreateTemp: %w", err)
	}

	_, fileName := filepath.Split(newFile.Name())
	currentDirectoryRecord := RecordEntity{
		DataFileName: fileName,
		Path:         pathIdentifier,
	}
	return currentDirectoryRecord, nil
}
