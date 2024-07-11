/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DanWlker/remind/internal/app/list"
	"github.com/DanWlker/remind/internal/config"
)

var allFlagList = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage:     "List all available todos",
	},
	Value: false,
}

var globalFlagList = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "List global todos",
	},
	Value: false,
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists todos",
	Long: `Lists todos, by default it attempts to list todos associated to
	this folder, use the -a flag to list all todos`,
	Run: func(cmd *cobra.Command, _ []string) {
		allFlag, err := cmd.Flags().GetBool(allFlagList.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: errGetBoolAllFlag: %w", err))
		}

		globalFlag, err := cmd.Flags().GetBool(globalFlagList.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: errGetBoolGlobalFlag: %w", err))
		}

		err = list.ListRun(allFlag, globalFlag)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("listRun: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(allFlagList.Name, allFlagList.Shorthand, allFlagList.Value, allFlagList.Usage)
	listCmd.Flags().BoolP(globalFlagList.Name, globalFlagList.Shorthand, globalFlagList.Value, globalFlagList.Usage)
}
