package remove

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/DanWlker/remind/internal/data"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

func removeTodoAssociatedWith(directory string, indexesToRemove map[int]bool) error {
	projectRecordEntity, err := record.GetRecordEntityWithIdentifier(directory)
	if err != nil {
		return fmt.Errorf("record.GetRecordEntityWithIdentifier: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	fullPath := filepath.Join(dataFolder, projectRecordEntity.DataFileName)
	todoList, err := data.GetTodoFromFile(fullPath)
	if err != nil {
		return fmt.Errorf("data.GetTodoFromFile: %w", err)
	}

	var newTodoList []data.TodoEntity
	for i, todo := range todoList {
		if _, shouldRemove := indexesToRemove[i]; shouldRemove {
			continue
		}
		newTodoList = append(newTodoList, todo)
	}

	if err := data.WriteTodoToFile(fullPath, newTodoList); err != nil {
		return fmt.Errorf("data.WriteTodoToFile: %w", err)
	}

	return nil
}

func removeAllTodosAssociatedWith(directory string) error {
	projectRecordEntity, err := record.GetRecordEntityWithIdentifier(directory)
	if err != nil {
		return fmt.Errorf("record.GetRecordEntityWithIdentifier: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	fullPath := filepath.Join(dataFolder, projectRecordEntity.DataFileName)
	if err := data.WriteTodoToFile(fullPath, []data.TodoEntity{}); err != nil {
		return fmt.Errorf("data.WriteTodoToFile: %w", err)
	}

	return nil
}

func RemoveRun(globalFlag, allFlag bool, args []string) error {
	indexesToRemove := make(map[int]bool)

	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("strconv.Atoi: %w", err)
		}

		indexesToRemove[i] = true
	}

	if globalFlag && allFlag {
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}
		if err := removeAllTodosAssociatedWith(homeRemoved); err != nil {
			return fmt.Errorf("removeAllTodosAssociatedWith: %w", err)
		}
		return nil
	}

	if globalFlag {
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}
		if err := removeTodoAssociatedWith(homeRemoved, indexesToRemove); err != nil {
			return fmt.Errorf("removeTodoAssociatedWith: %w", err)
		}

		return nil
	}

	homeRemovedProgramDir, err := shared.GetHomeRemovedWorkingDir()
	if err != nil {
		return fmt.Errorf("shared.GetHomeRemovedWorkingDir: %w", err)
	}

	if allFlag {
		if err := removeAllTodosAssociatedWith(homeRemovedProgramDir); err != nil {
			return fmt.Errorf("removeAllTodosAssociatedWith: %w", err)
		}
		return nil
	}

	if err := removeTodoAssociatedWith(homeRemovedProgramDir, indexesToRemove); err != nil {
		return fmt.Errorf("removeTodoAssociatedWith: %w", err)
	}

	return nil
}
