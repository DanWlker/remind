package list

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/DanWlker/remind/internal/data"
	i_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

func listOne(pathToFind string) error {
	projectRecordEntity, err := record.GetRecordEntityWithIdentifier(pathToFind)
	var errRecordDoesNotExist i_error.RecordDoesNotExistError
	if errors.As(err, &errRecordDoesNotExist) {
		recordIdentifier := errRecordDoesNotExist.ID
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}
		if recordIdentifier == homeRemoved {
			recordIdentifier = "$HOME"
		}
		fmt.Println("No record linked to this folder found: " + recordIdentifier)
		return nil
	} else if err != nil {
		return fmt.Errorf("record.GetRecordEntityWithIdentifier: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	err = data.PrettyPrintFile(
		filepath.Join(dataFolder, projectRecordEntity.DataFileName),
		func(todo string, index int) (string, error) {
			return fmt.Sprintf("\t%v. %v", index, todo), nil
		})
	if err != nil {
		return fmt.Errorf("data.PrettyPrintFile: %w", err)
	}
	return nil
}

func listConcurrently(item record.RecordEntity, dataFolder string) (<-chan string, error) {
	c := make(chan string)

	var header string
	homeRemoved, err := shared.GetHomeRemovedHomeDir()
	if err != nil {
		return nil, fmt.Errorf("shared.FormatRemoveHome: %w", err)
	}

	if item.Path == homeRemoved {
		header = "Global:\n"
	} else {
		header = item.Path + ":\n"
	}

	go func() {
		defer close(c)
		result, err := data.SPrettyPrintFile(
			filepath.Join(dataFolder, item.DataFileName),
			func(todo string, index int) (string, error) {
				return fmt.Sprintf("\t%v. %v", index, todo), nil
			},
		)
		if err != nil {
			c <- fmt.Sprintf("Error: Something went wrong: data.SPrettyPrintDataFile: %v", err)
			return
		}
		c <- header + result
	}()

	return c, nil
}

func listAll() error {
	items, err := record.GetFileContents()
	if err != nil {
		return fmt.Errorf("record.GetFileContents: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	var channelList []<-chan string
	for _, item := range items {
		c, err := listConcurrently(item, dataFolder)
		if err != nil {
			return fmt.Errorf("listConcurrently: %w", err)
		}

		channelList = append(channelList, c)
	}

	for _, channel := range channelList {
		fmt.Println(<-channel)
	}

	return nil
}

func ListRun(allFlag, globalFlag bool) error {
	// Check should list all
	if allFlag {
		if err := listAll(); err != nil {
			return fmt.Errorf("listAll: %w", err)
		}
		return nil
	}

	// Check should list global
	if globalFlag {
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}
		if err := listOne(homeRemoved); err != nil {
			return fmt.Errorf("globalFlag: listOne: %w", err)
		}
		return nil
	}

	// Attempt to get current directory and list reminders associated with it
	pathToFind, err := shared.GetHomeRemovedWorkingDir()
	var notUnderHomeError i_error.NotUnderHomeError
	if errors.As(err, &notUnderHomeError) {
		log.Println(
			notUnderHomeError.Error(),
		)
	} else if err != nil {
		return fmt.Errorf("shared.GetHomeRemovedWorkingDir: %w", err)
	}

	if err := listOne(pathToFind); err != nil {
		return fmt.Errorf("listOne: %v: %w", pathToFind, err)
	}
	return nil
}
