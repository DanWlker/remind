package edit

import (
	"fmt"
	"os"
	"slices"
	"strings"

	r_error "github.com/DanWlker/remind/internal/error"
	"github.com/DanWlker/remind/internal/pkg/data"
	"github.com/DanWlker/remind/internal/pkg/record"
	"github.com/DanWlker/remind/internal/pkg/shared"
)

func editTodoAssociatedWith(directory string) error {
	recordItems, errGetRecordFileContents := record.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	idx := slices.IndexFunc(recordItems, func(item record.RecordEntity) bool {
		return item.Path == directory
	})
	if idx == -1 {
		return &r_error.RecordDoesNotExistError{}
	}
	currentDirectoryRecord := &recordItems[idx]

	dataFolder, errGetDataFolder := data.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	dataFileFullPath := dataFolder + string(os.PathSeparator) + currentDirectoryRecord.DataFileName

	prettyPrintedString, errSPrettyPrintDataFile := data.SPrettyPrintDataFile(dataFileFullPath, nil)
	if errSPrettyPrintDataFile != nil {
		return fmt.Errorf("data.SPrettyPrintDataFile: %w", errSPrettyPrintDataFile)
	}

	result, errOpenDefaultEditor := shared.OpenDefaultEditor([]byte(prettyPrintedString))
	if errOpenDefaultEditor != nil {
		return fmt.Errorf("shared.OpenDefaultEditor: %w", errOpenDefaultEditor)
	}

	var todoList []data.TodoEntity
	for _, item := range strings.Split(strings.ReplaceAll(string(result), "\r\n", "\n"), "\n") {
		if item != "" {
			todoList = append(todoList, data.TodoEntity{Text: item})
		}
	}

	errWriteTodoToFile := data.WriteTodoToFile(dataFileFullPath, todoList)
	if errWriteTodoToFile != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", errWriteTodoToFile)
	}

	return nil
}

func EditRun(globalFlag bool) error {
	if globalFlag {
		errAddTodoAndAssociateTo := editTodoAssociatedWith("")
		if errAddTodoAndAssociateTo != nil {
			return fmt.Errorf("editTodoAssociatedWith: %w", errAddTodoAndAssociateTo)
		}
		return nil
	}

	homeRemCurrProExDir, errHomeRemCurrProExDir := shared.GetHomeRemovedCurrentProgramExecutionDirectory()
	if errHomeRemCurrProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errHomeRemCurrProExDir)
	}

	if errEditTodoAssociatedWith := editTodoAssociatedWith(homeRemCurrProExDir); errEditTodoAssociatedWith != nil {
		return fmt.Errorf("editTodoAssociatedWith: %w", errEditTodoAssociatedWith)
	}
	return nil
}
