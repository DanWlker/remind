/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/DanWlker/remind/internal/app/edit"
	"github.com/DanWlker/remind/internal/config"
	"github.com/spf13/cobra"
)

var globalFlag_edit = config.BoolFlagEntity{
	FlagEntity: config.FlagEntity{
		Name:      "global",
		Shorthand: "g",
		Usage:     "Edit the global todo list",
	},
	Value: false,
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"ed"},
	Short:   "Edits todos",
	Long: `Edits todos, by default it will let you to edit todos associated
	to this folder, use the -g flag to edit the global todo`,
	Run: func(cmd *cobra.Command, args []string) {
		globalFlag, errGetBool_global := cmd.Flags().GetBool(globalFlag_edit.Name)
		if errGetBool_global != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool_global))
		}

		errEditRun := edit.EditRun(globalFlag)
		if errEditRun != nil {
			cobra.CheckErr(fmt.Errorf("editRun: %w", errEditRun))
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolP(globalFlag_edit.Name, globalFlag_edit.Shorthand, globalFlag_edit.Value, globalFlag_edit.Usage)
}
