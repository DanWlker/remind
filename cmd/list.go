/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/DanWlker/remind/entity"
	r_error "github.com/DanWlker/remind/error"
	"github.com/DanWlker/remind/helper"
	"github.com/spf13/cobra"
)

var allFlag = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage:     "List all available todos",
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
		errListRun := listRun(cmd, args)
		if errListRun != nil {
			cobra.CheckErr(fmt.Errorf("listRun: %w", errListRun))
		}
	},
}

func listOne(fileFullPath string) {
	// TODO: Implement this
	fmt.Println(fileFullPath)
	// helper.ReadFromFile(fileFullPath)
}

func listAll() error {
	items, errGetRecordFileContents := helper.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	for _, item := range items {
		listOne(item.Path)
	}

	return nil
}

func listRun(cmd *cobra.Command, _ []string) error {
	shouldListAll, errGetBool := cmd.Flags().GetBool(allFlag.Name)
	if errGetBool != nil {
		return fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool)
	}

	if shouldListAll {
		errListAll := listAll()
		if errListAll != nil {
			return fmt.Errorf("listAll: %w", errListAll)
		}
		return nil
	}

	// Attempt to get current directory and list reminders associated with it
	currentProgramRunRemovedHomePath, errGetHomeRemovedFilePath := helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	var filePathNotStartsWithHomeErr *r_error.FilePathNotStartsWithHome
	if errors.As(errGetHomeRemovedFilePath, &filePathNotStartsWithHomeErr) {
		log.Println(
			fmt.Sprintf("Current program executed in path that does not include $HOME(%v), listRun:", filePathNotStartsWithHomeErr.HomeStr),
		)
	} else if errGetHomeRemovedFilePath != nil {
		return fmt.Errorf("helper.GetHomeRemovedFilePath: %w", errGetHomeRemovedFilePath)
	}

	listOne(currentProgramRunRemovedHomePath)
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(allFlag.Name, allFlag.Shorthand, allFlag.Value, allFlag.Usage)
}
