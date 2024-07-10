package record

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/DanWlker/remind/internal/config"
	"github.com/DanWlker/remind/internal/data"
	i_error "github.com/DanWlker/remind/internal/error"
)

func GetFile() (string, error) {
	dataFolder, err := data.GetFolder()
	if err != nil {
		return "", fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	fullpath := filepath.Join(dataFolder, config.DefaultDataRecordFullFileName)

	if _, err := os.Stat(fullpath); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(fullpath); err != nil {
			return "", fmt.Errorf("os.Create: %w", err)
		}
		globalRecordEntity, err := CreateNewRecord("")
		if err != nil {
			return "", fmt.Errorf("CreateNewRecord: %w", err)
		}

		if err := SetFileContents([]RecordEntity{globalRecordEntity}); err != nil {
			return "", fmt.Errorf("SetRecordFileContents: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("os.Stat: %w", err)
	}

	return fullpath, nil
}

func GetFileContents() ([]RecordEntity, error) {
	recordFile, err := GetFile()
	if err != nil {
		return nil, fmt.Errorf("GetRecordFile: %w", err)
	}

	f, err := os.Open(recordFile)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	var (
		items []RecordEntity
		dec   = yaml.NewDecoder(f)
	)

	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("in yaml decoding: %w", err)
	}

	return items, nil
}

func SetFileContents(items []RecordEntity) error {
	recordFile, err := GetFile()
	if err != nil {
		return fmt.Errorf("GetRecordFile: %w", err)
	}

	f, err := os.OpenFile(recordFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("in yaml encoding: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("closing file: %w", err)
	}

	return nil
}

func GetRecordEntityWithIdentifier(homeRemovedPath string) (RecordEntity, error) {
	allRecords, err := GetFileContents()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("GetRecordFileContents: %w", err)
	}

	for _, record := range allRecords {
		if record.Path == homeRemovedPath {
			return record, nil
		}
	}

	return RecordEntity{}, i_error.RecordDoesNotExistError{ID: homeRemovedPath}
}

func CreateNewRecord(pathIdentifier string) (RecordEntity, error) {
	dataFolder, err := data.GetFolder()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("GetDataFolder: %w", err)
	}

	newFile, err := os.CreateTemp(dataFolder, "*"+config.DefaultDataFileExtension)
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
