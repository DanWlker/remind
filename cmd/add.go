/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	"github.com/DanWlker/remind/helper"
	"github.com/spf13/cobra"
)

var globalFlag = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Adds the todos to a global todo list",
	},
	Value: false,
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a todo with directory context",
	Long:  "Adds a todo with directory context, by default it will associate the todo with the local directory. Use -g to bind it to the global $HOME todo list",
	Run: func(cmd *cobra.Command, args []string) {
		errAddRun := addRun(cmd, args)
		if errAddRun != nil {
			cobra.CheckErr(fmt.Errorf("addRun: %w", errAddRun))
		}
	},
}

func addTodoAndAssociateTo(directory string, todoListString []string) error {
	// Find the record in the record file
	recordItems, errGetRecordFileContents := helper.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	idx := slices.IndexFunc(recordItems, func(item entity.ProjectRecordEntity) bool { return item.Path == directory })

	var currentDirectoryRecord *entity.ProjectRecordEntity
	if idx == -1 {
		dataFolder, errGetDataFolder := helper.GetDataFolder()
		if errGetDataFolder != nil {
			return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
		}

		newFile, errCreateTemp := os.CreateTemp(dataFolder, "*"+constant.DEFAULT_DATA_FILE_EXTENSION)
		if errCreateTemp != nil {
			return fmt.Errorf("os.CreateTemp: %w", errCreateTemp)
		}

		currentDirectoryRecord = &entity.ProjectRecordEntity{
			DataFileName: newFile.Name(),
			Path:         directory,
		}
		recordItems = append(recordItems, *currentDirectoryRecord)
		helper.SetRecordFileContents(recordItems)
	} else {
		currentDirectoryRecord = &recordItems[idx]
	}

	fmt.Println("===============")
	fmt.Println(currentDirectoryRecord.DataFileName)
	fmt.Println("===============")

	// var todoList []entity.TodoEntity
	// for _, item := range todoListString {
	// 	todoList = append(todoList, entity.TodoEntity{Text: item})
	// }
	//
	// yamlTodoList, errMarshal := yaml.Marshal(todoList)
	// if errMarshal != nil {
	// 	return fmt.Errorf("yaml.Marshal: %w", errMarshal)
	// }
	//
	// dataFileName := currentDirectoryRecord.DataFileName
	// os.WriteFile(dataFileName, yamlTodoList, 0644)

	return nil
}

func addRun(cmd *cobra.Command, args []string) error {
	shouldAddToGlobal, errGetBool := cmd.Flags().GetBool(globalFlag.Name)
	if errGetBool != nil {
		return fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool)
	}

	if shouldAddToGlobal {
		errAddTodoAndAssociateTo := addTodoAndAssociateTo("~", args)
		if errAddTodoAndAssociateTo != nil {
			return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
		}
		return nil
	}

	homeRemCurrProExDir, errHomeRemCurrProExDir := helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	if errHomeRemCurrProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errHomeRemCurrProExDir)
	}

	errAddTodoAndAssociateTo := addTodoAndAssociateTo(homeRemCurrProExDir, args)
	if errAddTodoAndAssociateTo != nil {
		return fmt.Errorf("addTodoAndAssociateTo: %w", errAddTodoAndAssociateTo)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(globalFlag.Name, globalFlag.Shorthand, globalFlag.Value, globalFlag.Usage)
}
