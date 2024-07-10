package edit

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/DanWlker/remind/internal/data"
	r_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

func editTodoAssociatedWith(directory string) error {
	recordItems, err := record.GetFileContents()
	if err != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", err)
	}

	idx := slices.IndexFunc(recordItems, func(item record.RecordEntity) bool {
		return item.Path == directory
	})
	if idx == -1 {
		return r_error.RecordDoesNotExistError{ID: directory}
	}
	current := &recordItems[idx]

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	fullpath := filepath.Join(dataFolder, current.DataFileName)

	prettyPrintedString, err := data.SPrettyPrintFile(fullpath, nil)
	if err != nil {
		return fmt.Errorf("data.SPrettyPrintDataFile: %w", err)
	}

	result, err := shared.OpenDefaultEditor([]byte(prettyPrintedString))
	if err != nil {
		return fmt.Errorf("shared.OpenDefaultEditor: %w", err)
	}

	var (
		todoList []data.TodoEntity
		sc       = bufio.NewScanner(bytes.NewReader(result))
	)
	for sc.Scan() {
		if item := sc.Text(); item != "" {
			todoList = append(todoList, data.TodoEntity{Text: item})
		}
	}

	if err := data.WriteTodoToFile(fullpath, todoList); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}

	return nil
}

func EditRun(globalFlag bool) error {
	if globalFlag {
		if err := editTodoAssociatedWith(""); err != nil {
			return fmt.Errorf("editTodoAssociatedWith: %w", err)
		}
		return nil
	}

	dir, err := shared.GetHomeRemovedWorkingDir()
	if err != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", err)
	}

	if err := editTodoAssociatedWith(dir); err != nil {
		return fmt.Errorf("editTodoAssociatedWith: %w", err)
	}
	return nil
}
