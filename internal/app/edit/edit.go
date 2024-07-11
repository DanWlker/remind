package edit

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/DanWlker/remind/internal/data"
	i_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

func editTodoAssociatedWith(directory string) error {
	current, err := record.GetRecordEntityWithIdentifier(directory)
	var errRecordDoesNotExist i_error.RecordDoesNotExistError
	if errors.As(err, &errRecordDoesNotExist) {
		recordIdentifier := errRecordDoesNotExist.ID
		fmt.Println("No record linked to this folder found: " + recordIdentifier)
		return nil
	} else if err != nil {
		return fmt.Errorf("record.GetRecordEntityWithIdentifier: %w", err)
	}

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
		homeRemoved, err := shared.GetHomeRemovedHomeDir()
		if err != nil {
			return fmt.Errorf("shared.GetHomeRemovedHomeDir: %w", err)
		}
		if err := editTodoAssociatedWith(homeRemoved); err != nil {
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
