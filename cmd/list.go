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

var allFlag_list = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage:     "List all available todos",
	},
	Value: false,
}

var globalFlag_list = config.BoolFlagEntity{
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
	Run: func(cmd *cobra.Command, args []string) {
		allFlag, errGetBoolAllFlag := cmd.Flags().GetBool(allFlag_list.Name)
		if errGetBoolAllFlag != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: errGetBoolAllFlag: %w", errGetBoolAllFlag))
		}

		globalFlag, errGetBoolGlobalFlag := cmd.Flags().GetBool(globalFlag_list.Name)
		if errGetBoolGlobalFlag != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: errGetBoolGlobalFlag: %w", errGetBoolGlobalFlag))
		}
		errListRun := list.ListRun(allFlag, globalFlag)
		if errListRun != nil {
			cobra.CheckErr(fmt.Errorf("listRun: %w", errListRun))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(allFlag_list.Name, allFlag_list.Shorthand, allFlag_list.Value, allFlag_list.Usage)
	listCmd.Flags().BoolP(globalFlag_list.Name, globalFlag_list.Shorthand, globalFlag_list.Value, globalFlag_list.Usage)
}
