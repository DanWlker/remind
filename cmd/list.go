/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DanWlker/remind/helper"
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
		dataFolder, errGetDataFolder := helper.GetDataFolder()
		if errGetDataFolder != nil {
			// TODO: Change this when we optimize os.Stat from GetDataFolder
			cobra.CheckErr(fmt.Errorf("helper.GetDataFolder: %w", errGetDataFolder))
		}

		listAll, errGetBool := cmd.Flags().GetBool(allFlag)
		if errGetBool != nil {
			cobra.CheckErr(fmt.Errorf("cmd.Flags().GetBool: %w", errGetBool))
		}

		if listAll {
			items, errReadDir := os.ReadDir(dataFolder)
			if errReadDir != nil {
				cobra.CheckErr(fmt.Errorf("os.ReadDir(%w): %w", dataFolder, errReadDir))
			}
			for _, item := range items {
				if item.IsDir() {
					continue
				}

			}
		}

		items, err := ReadItems(dataFile)

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Println(items)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(allFlag, "a", false, "List all available todos")
}
