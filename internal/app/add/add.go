package add

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/DanWlker/remind/internal/pkg/data"
	"github.com/DanWlker/remind/internal/pkg/record"
	"github.com/DanWlker/remind/internal/pkg/shared"
)

func AddRun(globalFlag bool, args []string) error {
	if globalFlag {
		err := addTodoAndAssociateTo("", args)
		if err != nil {
			return fmt.Errorf("addTodoAndAssociateTo: %w", err)
		}
		return nil
	}

	dir, err := shared.GetHomeRemovedWorkingDir()
	if err != nil {
		return fmt.Errorf("shared.GetHomeRemovedWorkingDir: %w", err)
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
		return fmt.Errorf("record.GetFileContents: %w", err)
	}

	idx := slices.IndexFunc(recordItems, func(item record.RecordEntity) bool {
		return item.Path == directory
	})

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	var current *record.RecordEntity
	if idx == -1 {
		tmp, err := record.CreateNewRecord(directory)
		if err != nil {
			return fmt.Errorf("record.CreateNewRecord: %w", err)
		}
		current = &tmp
		recordItems = append(recordItems, tmp)
		if err := record.SetFileContents(recordItems); err != nil {
			return fmt.Errorf("record.SetFileContents: %w", err)
		}
	} else {
		current = &recordItems[idx]
	}

	// Read the file, it will exist if it reaches here
	fullPath := filepath.Join(dataFolder, current.DataFileName)

	todoList, err := data.GetTodoFromFile(fullPath)
	if err != nil {
		return fmt.Errorf("data.GetTodoFromFile: %w", err)
	}

	for _, item := range todoListString {
		todoList = append(todoList, data.TodoEntity{Text: item})
	}

	err = data.WriteTodoToFile(fullPath, todoList)
	if err != nil {
		return fmt.Errorf("data.WriteTodoToFile: %w", err)
	}
	return nil
}
