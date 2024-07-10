/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DanWlker/remind/internal/app/remove"
	"github.com/DanWlker/remind/internal/config"
)

var globalFlagRemove = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Removes the todos in the global list",
	},
	Value: false,
}

var allFlagRemove = config.BoolFlagEntity{
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
		globalFlag, err := cmd.Flags().GetBool(globalFlagRemove.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		allFlag, err := cmd.Flags().GetBool(allFlagRemove.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		if err := remove.RemoveRun(globalFlag, allFlag, args); err != nil {
			cobra.CheckErr(fmt.Errorf("removeRun: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP(globalFlagRemove.Name, globalFlagRemove.Shorthand, globalFlagRemove.Value, globalFlagRemove.Usage)
	removeCmd.Flags().BoolP(allFlagRemove.Name, allFlagRemove.Shorthand, allFlagRemove.Value, allFlagRemove.Usage)
}
