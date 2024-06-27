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
	projectRecordEntity, errFindProjectRecordEntity := record.GetProjectRecordFromFileWith(pathToFind)
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

	dataFolder, errGetDataFolder := data.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	if errPrettyPrintFile := data.PrettyPrintDataFile(dataFolder+string(os.PathSeparator)+projectRecordEntity.DataFileName, "  "); errPrettyPrintFile != nil {
		return fmt.Errorf("helper.PrettyPrintDataFile: %w", errPrettyPrintFile)
	}
	return nil
}

func listAll() error {
	items, errGetRecordFileContents := record.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	dataFolder, errGetDataFolder := data.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	for _, item := range items {
		if item.Path == "" {
			fmt.Println("Global:")
		} else {
			fmt.Println(item.Path + ":")
		}

		if errPrettyPrintDataFile := data.PrettyPrintDataFile(dataFolder+string(os.PathSeparator)+item.DataFileName, "  "); errPrettyPrintDataFile != nil {
			return errPrettyPrintDataFile
		}

		fmt.Println("")
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
	pathToFind, errGetHomeRemovedFilePath := shared.GetHomeRemovedCurrentProgramExecutionDirectory()
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
