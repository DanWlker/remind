package add

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/DanWlker/remind/internal/data"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

func AddRun(globalFlag bool, args []string) error {
	if globalFlag {
		if err := addTodoAndAssociateTo("", args); err != nil {
			return fmt.Errorf("addTodoAndAssociateTo: %w", err)
		}
		return nil
	}

	dir, err := shared.GetHomeRemovedWorkingDir()
	if err != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", err)
	}

	if err := addTodoAndAssociateTo(dir, args); err != nil {
		return fmt.Errorf("addTodoAndAssociateTo: %w", err)
	}

	return nil
}

func addTodoAndAssociateTo(directory string, todoListString []string) error {
	// Find the record in the record file
	recordItems, err := record.GetFileContents()
	if err != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", err)
	}

	idx := slices.IndexFunc(recordItems, func(item record.RecordEntity) bool {
		return item.Path == directory
	})

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	var current *record.RecordEntity
	if idx == -1 {
		// Record was not found.
		// Create a new one and add it to recordItems.

		tmp, err := record.CreateNewRecord(directory)
		if err != nil {
			return fmt.Errorf("helper.CreateNewRecord: %w", err)
		}
		current = &tmp
		recordItems = append(recordItems, tmp)
		if err := record.SetFileContents(recordItems); err != nil {
			return fmt.Errorf("helper.SetRecordFileContents: %w", err)
		}
	} else {
		current = &recordItems[idx]
	}

	// Read the file, it will exist if it reaches here
	fullpath := filepath.Join(dataFolder, current.DataFileName)
	todoList, err := data.GetTodoFromFile(fullpath)
	if err != nil {
		return fmt.Errorf("helper.ReadFromFile: %w", err)
	}

	for _, item := range todoListString {
		todoList = append(todoList, data.TodoEntity{Text: item})
	}

	if err := data.WriteTodoToFile(fullpath, todoList); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}
	return nil
}
