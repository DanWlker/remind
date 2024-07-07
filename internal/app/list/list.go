package list

import (
	"errors"
	"fmt"
	"log"
	"os"

	i_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/pkg/data"
	"github.com/DanWlker/remind/internal/pkg/record"
	"github.com/DanWlker/remind/internal/pkg/shared"
)

func listOne(pathToFind string) error {
	projectRecordEntity, errFindProjectRecordEntity := record.GetRecordEntityWithIdentifier(pathToFind)
	var errRecordDoesNotExist *i_error.RecordDoesNotExistError
	if errors.As(errFindProjectRecordEntity, &errRecordDoesNotExist) {
		recordIdentifier := errRecordDoesNotExist.RecordIdentifier
		if recordIdentifier == "" {
			recordIdentifier = "$HOME"
		}
		fmt.Println("No record linked to this folder found: " + recordIdentifier)
		return nil
	} else if errFindProjectRecordEntity != nil {
		return fmt.Errorf("helper.FindProjectRecordFromFileWith: %w", errFindProjectRecordEntity)
	}

	dataFolder, errGetDataFolder := data.GetFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	if errPrettyPrintFile := data.PrettyPrintFile(
		dataFolder+string(os.PathSeparator)+projectRecordEntity.DataFileName,
		func(todo string, index int) string {
			return fmt.Sprintf("\t%v. %v", index, todo)
		}); errPrettyPrintFile != nil {
		return fmt.Errorf("helper.PrettyPrintDataFile: %w", errPrettyPrintFile)
	}
	return nil
}

func listConcurrently(item record.RecordEntity, dataFolder string) (chan string, error) {
	var c chan string = make(chan string)

	var header string
	if item.Path == "" {
		header = "Global:\n"
	} else {
		header = item.Path + ":\n"
	}

	go func() {
		result, errPrettyPrintDataFile := data.SPrettyPrintFile(
			dataFolder+string(os.PathSeparator)+item.DataFileName,
			func(todo string, index int) string {
				return fmt.Sprintf("\t%v. %v", index, todo)
			},
		)
		if errPrettyPrintDataFile != nil {
			c <- fmt.Sprintf("Error: Something went wrong: data.SPrettyPrintDataFile: %v", errPrettyPrintDataFile)
		}
		c <- header + result
	}()

	return c, nil
}

func listAll() error {
	items, errGetRecordFileContents := record.GetFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	dataFolder, errGetDataFolder := data.GetFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	var channelList []chan string
	for _, item := range items {
		c, errlistConcurrently := listConcurrently(item, dataFolder)
		if errlistConcurrently != nil {
			log.Println("listConcurrently: %w", errlistConcurrently)
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
		if errListAll := listAll(); errListAll != nil {
			return fmt.Errorf("listAll: %w", errListAll)
		}
		return nil
	}

	// Check should list global
	if globalFlag {
		if errListOneGlobal := listOne(""); errListOneGlobal != nil {
			return fmt.Errorf("listOne: shouldListGlobal: %w", errListOneGlobal)
		}
		return nil
	}

	// Attempt to get current directory and list reminders associated with it
	pathToFind, errGetHomeRemovedFilePath := shared.GetHomeRemovedWorkingDir()
	var filePathNotStartsWithHomeErr *i_error.FilePathNotStartsWithHome
	if errors.As(errGetHomeRemovedFilePath, &filePathNotStartsWithHomeErr) {
		log.Println(
			filePathNotStartsWithHomeErr.Error(),
		)
	} else if errGetHomeRemovedFilePath != nil {
		return fmt.Errorf("helper.GetHomeRemovedFilePath: %w", errGetHomeRemovedFilePath)
	}
	if errListOneLocal := listOne(pathToFind); errListOneLocal != nil {
		return fmt.Errorf("listOne: %v: %w", pathToFind, errListOneLocal)
	}
	return nil
}
