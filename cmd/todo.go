/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"

	"github.com/DanWlker/remind/entity"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Adds a new todo tied to the current folder",
	Long: `Todo will write a new entry to a symbolically linked file in the
	current project. This file can also be found in $HOME/remind/. If the
	file does not exist, it will ask to provide a name for the file or
	symbolically link an existing file.`,
	Run: todoRun,
}

func todoRun(cmd *cobra.Command, args []string) {
	items := []entity.TodoEntity{}
	for _, x := range args {
		items = append(items, entity.TodoEntity{Text: x})
	}
	fmt.Printf("%#v\n", items)
}

func init() {
	rootCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
