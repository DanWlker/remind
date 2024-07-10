/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DanWlker/remind/internal/app/edit"
	"github.com/DanWlker/remind/internal/config"
)

var globalFlagEdit = config.BoolFlagEntity{
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
	Run: func(cmd *cobra.Command, _ []string) {
		globalFlag, err := cmd.Flags().GetBool(globalFlagEdit.Name)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", err))
		}

		err = edit.EditRun(globalFlag)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("editRun: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolP(globalFlagEdit.Name, globalFlagEdit.Shorthand, globalFlagEdit.Value, globalFlagEdit.Usage)
}
