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

var globalFlag_add = config.BoolFlagEntity{
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
		globalFlag, err := cmd.Flags().GetBool(globalFlag_add.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}
		if err := add.AddRun(globalFlag, args); err != nil {
			cobra.CheckErr(fmt.Errorf("addRun: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(globalFlag_add.Name, globalFlag_add.Shorthand, globalFlag_add.Value, globalFlag_add.Usage)
}
