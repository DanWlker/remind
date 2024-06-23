/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/DanWlker/remind/constant"
	"github.com/DanWlker/remind/entity"
	r_error "github.com/DanWlker/remind/error"
	"github.com/DanWlker/remind/helper"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

const allFlag = "all"

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
	dataFolder, errGetDataFolder := helper.GetDataFolder()
	if errGetDataFolder != nil {
		return fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder)
	}

	defaultDataRecordFileFullPath := dataFolder + constant.DEFAULT_DATA_RECORD_FULL_FILE_NAME

	_, errStat := os.Stat(defaultDataRecordFileFullPath)
	if errors.Is(errStat, os.ErrNotExist) {
		_, errCreate := os.Create(defaultDataRecordFileFullPath)
		if errCreate != nil {
			return fmt.Errorf("os.Create: %w", errStat)
		}
	} else if errStat != nil {
		return fmt.Errorf("os.Stat: %w", errStat)
	}

	recordFile, errReadFile := os.ReadFile(defaultDataRecordFileFullPath)
	if errReadFile != nil {
		return fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []entity.ProjectRecordEntity
	if errUnmarshal := yaml.Unmarshal(recordFile, &items); errUnmarshal != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	for _, item := range items {
		listOne(item.Path)
	}

	return nil
}

func listRun(cmd *cobra.Command, args []string) error {
	shouldListAll, errGetBool := cmd.Flags().GetBool(allFlag)
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
	currentProgramRunFullPath, errGetCurrentProgramExecutionDir := helper.GetCurrentProgramExecutionDirectory()
	if errGetCurrentProgramExecutionDir != nil {
		return fmt.Errorf("helper.GetCurrentProgramExecutionDirectory: %w", errGetCurrentProgramExecutionDir)
	}

	currentProgramRunRemovedHomePath, errGetHomeRemovedFilePath := helper.GetHomeRemovedPath(currentProgramRunFullPath)
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

	listCmd.Flags().BoolP(allFlag, "a", false, "List all available todos")
}
