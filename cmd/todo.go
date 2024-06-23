/*
Copyright © 2024 DanWlker
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/DanWlker/remind/entity"
	"github.com/goccy/go-yaml"
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
	// items, _ := ReadItems(helper.GetDataFile())
	// for _, x := range args {
	// 	items = append(items, entity.TodoEntity{Text: x})
	// }
	// err := SaveItems(helper.GetDataFile(), items)
	//
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// }
	//
	// fmt.Printf("%#v\n", items)
}

func SaveItems(filename string, items []entity.TodoEntity) error {
	b, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	fmt.Println(string(filename))

	err = os.WriteFile(filename, b, 0644)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ReadItems(filename string) ([]entity.TodoEntity, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		debug.PrintStack()
		log.Println(err)
		return []entity.TodoEntity{}, err
	}
	var items []entity.TodoEntity
	if err := yaml.Unmarshal(b, &items); err != nil {
		debug.PrintStack()
		log.Println(err)
		return []entity.TodoEntity{}, err
	}
	return items, nil
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
