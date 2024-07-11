/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DanWlker/remind/internal/app/add"
	"github.com/DanWlker/remind/internal/config"
)

var globalFlagAdd = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
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
	Long: `Adds a todo with directory context, by default it will
	associate the todo with the local directory. Use -g to bind it to the
	global $HOME todo list`,
	Run: func(cmd *cobra.Command, args []string) {
		globalFlag, err := cmd.Flags().GetBool(globalFlagAdd.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		err = add.AddRun(globalFlag, args)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("addRun: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(globalFlagAdd.Name, globalFlagAdd.Shorthand, globalFlagAdd.Value, globalFlagAdd.Usage)
}
