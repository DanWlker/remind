/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

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

func addTodoAndAssociateTo(directory string, todoList []string) {
	// TODO: Check file name associated with the current directory from record file

	// TODO: Create the file if it doesn't exist

	// TODO: Write to the file
}

func addRun(cmd *cobra.Command, args []string) error {
	shouldAddToGlobal, errGetBool := cmd.Flags().GetBool(globalFlag.Name)
	if errGetBool != nil {
		return fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool)
	}

	if shouldAddToGlobal {
		addTodoAndAssociateTo("~", args)
		return nil
	}

	homeRemCurrProExDir, errHomeRemCurrProExDir := helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	if errHomeRemCurrProExDir != nil {
		return fmt.Errorf("helper.GetHomeRemovedCurrentProgramExecutionDirectory: %w", errHomeRemCurrProExDir)
	}

	addTodoAndAssociateTo(homeRemCurrProExDir, args)

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(globalFlag.Name, globalFlag.Shorthand, globalFlag.Value, globalFlag.Usage)
}
