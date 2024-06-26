/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DanWlker/remind/entity"
	"github.com/DanWlker/remind/helper"
	"github.com/spf13/cobra"
)

var globalFlag_remove = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Removes the todos in the global list",
	},
	Value: false,
}

var allFlag_remove = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage:     "Removes all the todos for the chosen directory. By default the chosen directory is the local directory, use -g to switch to remove from the global list",
	},
	Value: false,
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes todos with directory context",
	Long:    "Removes todos with directory context, by default it will attempt to remove todos associated with the local directory. Use -g to refer to the global $HOME todo list",
	Run: func(cmd *cobra.Command, args []string) {

		globalFlag, errGetBool_global := cmd.Flags().GetBool(globalFlag_remove.Name)
		if errGetBool_global != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool_global))
		}

		allFlag, errGetBool_all := cmd.Flags().GetBool(allFlag_remove.Name)
		if errGetBool_all != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool_all))
		}

		if errRemoveRun := removeRun(globalFlag, allFlag, args); errRemoveRun != nil {
			cobra.CheckErr(fmt.Errorf("removeRun: %w", errGetBool_all))
		}
	},
}

func removeTodoAssociatedWith(directory string, indexesToRemove map[int]bool) error {
	projectRecordEntity, errGetProjectRecordFromFileWith := helper.GetProjectRecordFromFileWith(directory)
	if errGetProjectRecordFromFileWith != nil {
		return fmt.Errorf("helper.GetProjectRecordFromFileWith: %w", errGetProjectRecordFromFileWith)
	}

	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	todoList, errGetTodoFromDataFile := helper.GetTodoFromDataFile(dataFolder + string(os.PathSeparator) + projectRecordEntity.DataFileName)
	if errGetTodoFromDataFile != nil {
		return fmt.Errorf("helper.GetTodoFromDataFile: %w", errGetTodoFromDataFile)
	}

	var newTodoList []entity.TodoEntity
	for i, todo := range todoList {
		_, shouldRemove := indexesToRemove[i]
		if !shouldRemove {
			continue
		}
		newTodoList = append(newTodoList, todo)
	}

	dataFileFullPath := dataFolder + string(os.PathSeparator) + projectRecordEntity.DataFileName
	if err := helper.WriteTodoToFile(dataFileFullPath, newTodoList); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}

	return nil
}

func removeAllTodosAssociatedWith(directory string) error {
	projectRecordEntity, errGetProjectRecordFromFileWith := helper.GetProjectRecordFromFileWith(directory)
	if errGetProjectRecordFromFileWith != nil {
		return fmt.Errorf("helper.GetProjectRecordFromFileWith: %w", errGetProjectRecordFromFileWith)
	}

	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	dataFileFullPath := dataFolder + string(os.PathSeparator) + projectRecordEntity.DataFileName
	if err := helper.WriteTodoToFile(dataFileFullPath, []entity.TodoEntity{}); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}

	return nil
}

func removeRun(globalFlag, allFlag bool, args []string) error {
	indexesToRemove := make(map[int]bool)

	for _, arg := range args {
		i, errAtoi := strconv.Atoi(arg)
		if errAtoi != nil {
			return fmt.Errorf("strconv.Atoi: %w", errAtoi)
		}

		indexesToRemove[i] = true
	}

	if globalFlag && allFlag {
		if err := removeAllTodosAssociatedWith(""); err != nil {
			return fmt.Errorf("removeAllTodosAssociatedWith: %w", err)
		}
		return nil
	}
	if globalFlag {
		if err := removeTodoAssociatedWith("", indexesToRemove); err != nil {
			return fmt.Errorf("removeTodoAssociatedWith: %w", err)
		}

		return nil
	}

	homeRemovedProgramDir, errGetHomeRemProExDir := helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	if errGetHomeRemProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errGetHomeRemProExDir)
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

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP(globalFlag_remove.Name, globalFlag_remove.Shorthand, globalFlag_remove.Value, globalFlag_remove.Usage)
	removeCmd.Flags().BoolP(allFlag_remove.Name, allFlag_remove.Shorthand, allFlag_remove.Value, allFlag_remove.Usage)
}
