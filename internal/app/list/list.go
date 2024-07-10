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
		if recordIdentifier == "" {
			recordIdentifier = "$HOME"
		}
		fmt.Println("No record linked to this folder found: " + recordIdentifier)
		return nil
	}
	if err != nil {
		return fmt.Errorf("helper.FindProjectRecordFromFileWith: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	err = data.PrettyPrintFile(
		filepath.Join(dataFolder, projectRecordEntity.DataFileName),
		func(todo string, index int) string {
			return fmt.Sprintf("\t%v. %v", index, todo)
		},
	)
	if err != nil {
		return fmt.Errorf("helper.PrettyPrintDataFile: %w", err)
	}
	return nil
}

func listConcurrently(item record.RecordEntity, dataFolder string) (<-chan string, error) {
	c := make(chan string)

	var header string
	if item.Path == "" {
		header = "Global:\n"
	} else {
		header = item.Path + ":\n"
	}

	go func() {
		defer close(c)

		result, err := data.SPrettyPrintFile(
			filepath.Join(dataFolder, item.DataFileName),
			func(todo string, index int) string {
				return fmt.Sprintf("\t%v. %v", index, todo)
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
		return fmt.Errorf("helper.GetRecordFileContents: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	var channelList []<-chan string
	for _, item := range items {
		c, err := listConcurrently(item, dataFolder)
		if err != nil {
			log.Println("listConcurrently: %w", err)
			continue
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
		if err := listOne(""); err != nil {
			return fmt.Errorf("listOne: shouldListGlobal: %w", err)
		}
		return nil
	}

	// Attempt to get current directory and list reminders associated with it
	pathToFind, err := shared.GetHomeRemovedWorkingDir()

	var notUnderHomeErr i_error.NotUnderHomeError
	if errors.As(err, &notUnderHomeErr) {
		log.Println(notUnderHomeErr.Error())
	} else if err != nil {
		return fmt.Errorf("helper.GetHomeRemovedFilePath: %w", err)
	}
	if err := listOne(pathToFind); err != nil {
		return fmt.Errorf("listOne: %v: %w", pathToFind, err)
	}
	return nil
}
