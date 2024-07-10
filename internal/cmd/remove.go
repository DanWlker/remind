/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/DanWlker/remind/internal/config"
	"github.com/DanWlker/remind/internal/data"
	"github.com/DanWlker/remind/internal/record"
	"github.com/DanWlker/remind/internal/shared"
)

var globalFlag_remove = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Removes the todos in the global list",
	},
	Value: false,
}

var allFlag_remove = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage: `Removes all the todos for the chosen directory. By
		default the chosen directory is the local directory, use -g to
		remove from the global list`,
	},
	Value: false,
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes todos with directory context",
	Long: `Removes todos with directory context, by default it will attempt
	to remove todos associated with the local directory. Use -g to refer to
	the global $HOME todo list`,
	Run: func(cmd *cobra.Command, args []string) {

		globalFlag, err := cmd.Flags().GetBool(globalFlag_remove.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		allFlag, err := cmd.Flags().GetBool(allFlag_remove.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		if err := removeRun(globalFlag, allFlag, args); err != nil {
			cobra.CheckErr(fmt.Errorf("removeRun: %w", err))
		}
	},
}

func removeTodoAssociatedWith(directory string, indexesToRemove map[int]bool) error {
	projectRecordEntity, err := record.GetRecordEntityWithIdentifier(directory)
	if err != nil {
		return fmt.Errorf("helper.GetProjectRecordFromFileWith: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	fullpath := filepath.Join(dataFolder, projectRecordEntity.DataFileName)

	todoList, err := data.GetTodoFromFile(fullpath)
	if err != nil {
		return fmt.Errorf("helper.GetTodoFromDataFile: %w", err)
	}

	var newTodoList []data.TodoEntity
	for i, todo := range todoList {
		if _, shouldRemove := indexesToRemove[i]; shouldRemove {
			continue
		}
		newTodoList = append(newTodoList, todo)
	}

	if err := data.WriteTodoToFile(fullpath, newTodoList); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}

	return nil
}

func removeAllTodosAssociatedWith(directory string) error {
	projectRecordEntity, err := record.GetRecordEntityWithIdentifier(directory)
	if err != nil {
		return fmt.Errorf("helper.GetProjectRecordFromFileWith: %w", err)
	}

	dataFolder, err := data.GetFolder()
	if err != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", err)
	}

	dataFileFullPath := filepath.Join(dataFolder, projectRecordEntity.DataFileName)
	if err := data.WriteTodoToFile(dataFileFullPath, []data.TodoEntity{}); err != nil {
		return fmt.Errorf("helper.WriteTodoToFile: %w", err)
	}

	return nil
}

func removeRun(globalFlag, allFlag bool, args []string) error {
	indexesToRemove := make(map[int]bool)

	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("strconv.Atoi: %w", err)
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

	homeRemovedProgramDir, err := shared.GetHomeRemovedWorkingDir()
	if err != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", err)
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
