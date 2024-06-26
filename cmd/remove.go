/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"
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

		fmt.Println("remove called")
		if errRemoveRun := removeRun(globalFlag, allFlag, args); errRemoveRun != nil {
			cobra.CheckErr(fmt.Errorf("removeRun: %w", errGetBool_all))
		}
	},
}

func removeTodoAssociatedWith(directory string, indexesToRemove []int) error {
	return nil
}

func removeAllTodosAssociatedWith(directory string) error {
	return nil
}

func removeRun(globalFlag, allFlag bool, args []string) error {
	var indexesToRemove []int
	for _, arg := range args {
		i, errAtoi := strconv.Atoi(arg)
		if errAtoi != nil {
			return fmt.Errorf("strconv.Atoi: %w", errAtoi)
		}
		indexesToRemove = append(indexesToRemove, i)
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
