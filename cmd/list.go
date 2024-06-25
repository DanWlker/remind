/*
list
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/DanWlker/remind/entity"
	r_error "github.com/DanWlker/remind/error"
	"github.com/DanWlker/remind/helper"
	"github.com/spf13/cobra"
)

var allFlag_list = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
		Name:      "all",
		Shorthand: "a",
		Usage:     "List all available todos",
	},
	Value: false,
}

var globalFlag_list = entity.BoolFlagEntity{
	FlagEntity: entity.FlagEntity{
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
		errListRun := listRun(cmd, args)
		if errListRun != nil {
			cobra.CheckErr(fmt.Errorf("listRun: %w", errListRun))
		}
	},
}

func listOne(pathToFind string) error {
	projectRecordEntity, errFindProjectRecordEntity := helper.FindProjectRecordFromFileWith(pathToFind)
	var errRecordDoesNotExist *r_error.RecordDoesNotExistError
	if errors.As(errFindProjectRecordEntity, &errRecordDoesNotExist) {
		recordIdentifier := errRecordDoesNotExist.RecordIdentifier
		if recordIdentifier == "" {
			recordIdentifier = "$HOME"
		}
		fmt.Println("No record linked to this folder found: " + recordIdentifier)
		return nil
	} else if errFindProjectRecordEntity != nil {
		return fmt.Errorf("helper.FindProjectRecordFromFileWith: %w", errFindProjectRecordEntity)
	}

	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	if errPrettyPrintFile := helper.PrettyPrintDataFile(dataFolder+string(os.PathSeparator)+projectRecordEntity.DataFileName, "  "); errPrettyPrintFile != nil {
		return fmt.Errorf("helper.PrettyPrintDataFile: %w", errPrettyPrintFile)
	}
	return nil
}

func listAll() error {
	items, errGetRecordFileContents := helper.GetRecordFileContents()
	if errGetRecordFileContents != nil {
		return fmt.Errorf("helper.GetRecordFileContents: %w", errGetRecordFileContents)
	}

	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	for _, item := range items {
		if item.Path == "" {
			fmt.Println("Global:")
		} else {
			fmt.Println(item.Path + ":")
		}

		if errPrettyPrintDataFile := helper.PrettyPrintDataFile(dataFolder+string(os.PathSeparator)+item.DataFileName, "  "); errPrettyPrintDataFile != nil {
			return errPrettyPrintDataFile
		}

		fmt.Println("")
	}

	return nil
}

func listRun(cmd *cobra.Command, _ []string) error {
	// Check should list all
	shouldListAll, errGetBoolAllFlag := cmd.Flags().GetBool(allFlag_list.Name)
	if errGetBoolAllFlag != nil {
		return fmt.Errorf("cmd.Flags().GetBool: errGetBoolAllFlag: %w", errGetBoolAllFlag)
	}
	if shouldListAll {
		if errListAll := listAll(); errListAll != nil {
			return fmt.Errorf("listAll: %w", errListAll)
		}
		return nil
	}

	// Check should list global
	shoulListGlobal, errGetBoolGlobalFlag := cmd.Flags().GetBool(globalFlag_list.Name)
	if errGetBoolGlobalFlag != nil {
		return fmt.Errorf("cmd.Flags().GetBool: errGetBoolGlobalFlag: %w", errGetBoolGlobalFlag)
	}
	var pathToFind string
	if shoulListGlobal {
		if errListOneGlobal := listOne(""); errListOneGlobal != nil {
			return fmt.Errorf("listOne: shouldListGlobal: %w", errListOneGlobal)
		}
		return nil
	}

	// Attempt to get current directory and list reminders associated with it
	var errGetHomeRemovedFilePath error
	pathToFind, errGetHomeRemovedFilePath = helper.GetHomeRemovedCurrentProgramExecutionDirectory()
	var filePathNotStartsWithHomeErr *r_error.FilePathNotStartsWithHome
	if errors.As(errGetHomeRemovedFilePath, &filePathNotStartsWithHomeErr) {
		log.Println(
			filePathNotStartsWithHomeErr.Error(),
		)
	} else if errGetHomeRemovedFilePath != nil {
		return fmt.Errorf("helper.GetHomeRemovedFilePath: %w", errGetHomeRemovedFilePath)
	}
	if errListOneLocal := listOne(pathToFind); errListOneLocal != nil {
		return fmt.Errorf("listOne: %v: %w", pathToFind, errListOneLocal)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(allFlag_list.Name, allFlag_list.Shorthand, allFlag_list.Value, allFlag_list.Usage)
	listCmd.Flags().BoolP(globalFlag_list.Name, globalFlag_list.Shorthand, globalFlag_list.Value, globalFlag_list.Usage)
}
