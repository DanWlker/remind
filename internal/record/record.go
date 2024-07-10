package record

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/DanWlker/remind/internal/config"
	"github.com/DanWlker/remind/internal/data"
	i_error "github.com/DanWlker/remind/internal/error"
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

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("dec.Decode: %w", err)
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

	enc := yaml.NewEncoder(f)
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("enc.Encode: %w", err)
	}

	return nil
}

func GetRecordEntityWithIdentifier(homeRemovedPath string) (RecordEntity, error) {
	allRecords, err := GetFileContents()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("GetFileContents: %w", err)
	}

	for _, record := range allRecords {
		if record.Path == homeRemovedPath {
			return record, nil
		}
	}

	return RecordEntity{}, i_error.RecordDoesNotExistError{
		ID: homeRemovedPath,
	}
}

// TODO: This could leak temp files, perhaps need an api that returns the record entity but create the record for the path if it doesn't exist. TLDR instead of exposing a "Create" method directly, expose only the function to get the record
func CreateNewRecord(pathIdentifier string) (RecordEntity, error) {
	dataFolder, err := data.GetFolder()
	if err != nil {
		return RecordEntity{}, fmt.Errorf("data.GetFolder: %w", err)
	}

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
