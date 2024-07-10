package edit

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"slices"

	r_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/pkg/data"
	"github.com/DanWlker/remind/internal/pkg/record"
	"github.com/DanWlker/remind/internal/pkg/shared"
)

func editTodoAssociatedWith(directory string) error {
	recordItems, err := record.GetFileContents()
	if err != nil {
		return fmt.Errorf("record.GetFileContents: %w", err)
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
		return fmt.Errorf("data.GetFolder: %w", err)
	}

	fullPath := filepath.Join(dataFolder, current.DataFileName)

	prettyPrintedString, err := data.SPrettyPrintFile(fullPath, nil)
	if err != nil {
		return fmt.Errorf("data.SPrettyPrintDataFile: %w", err)
	}

	result, err := shared.OpenDefaultEditor([]byte(prettyPrintedString))
	if err != nil {
		return fmt.Errorf("shared.OpenDefaultEditor: %w", err)
	}

	var (
		todoList []data.TodoEntity
		// A bufio.Scanner can read input (from a file, a byte slice, a
		// string, etc) line by line, implicitly converting \r\n into
		// \n.
		sc = bufio.NewScanner(bytes.NewReader(result))
	)

	for sc.Scan() {
		if item := sc.Text(); item != "" {
			todoList = append(todoList, data.TodoEntity{Text: item})
		}
	}

	if err := data.WriteTodoToFile(fullPath, todoList); err != nil {
		return fmt.Errorf("data.WriteTodoToFile: %w", err)
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
		return fmt.Errorf("shared.GetHomeRemovedWorkingDir: %w", err)
	}

	if err := editTodoAssociatedWith(dir); err != nil {
		return fmt.Errorf("editTodoAssociatedWith: %w", err)
	}
	return nil
}
