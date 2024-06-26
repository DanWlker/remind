package data

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/DanWlker/remind/internal/config"
	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

// This does not create the file if it doesn't exist
func GetTodoFromDataFile(dataFileFullPath string) ([]TodoEntity, error) {
	file, errReadFile := os.ReadFile(dataFileFullPath)
	if errReadFile != nil {
		return []TodoEntity{}, fmt.Errorf("os.ReadFile: %w", errReadFile)
	}

	var items []TodoEntity
	if errUnmarshal := yaml.Unmarshal(file, &items); errUnmarshal != nil {
		return []TodoEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil
}

func WriteTodoToFile(dataFileFullPath string, todoList []TodoEntity) error {
	yamlTodoList, errMarshal := yaml.Marshal(todoList)
	if errMarshal != nil {
		return fmt.Errorf("yaml.Marshal: %w", errMarshal)
	}

	errWriteFile := os.WriteFile(dataFileFullPath, yamlTodoList, 0644)
	if errWriteFile != nil {
		return fmt.Errorf("os.WriteFile: %w", errWriteFile)
	}
	return nil
}

func GetDataFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(config.USER_DEFINED_DATA_FOLDER))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + config.DEFAULT_DATA_PATH_AFTER_HOME
	}

	if errMkDirAll := os.MkdirAll(dataFolder, 0770); errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return dataFolder, nil
}

func SPrettyPrintDataFile(dataFileFullPath string, editText func(todo string, index int) string) (string, error) {
	var b bytes.Buffer
	todoList, errGetTodoFromDataFile := GetTodoFromDataFile(dataFileFullPath)
	if errGetTodoFromDataFile != nil {
		return "", fmt.Errorf("GetTodoFromDataFile: %w", errGetTodoFromDataFile)
	}

	for i, todo := range todoList {
		if editText == nil {
			b.WriteString(todo.Text)
		} else {
			b.WriteString(editText(todo.Text, i))
		}
		b.WriteString("\n")
	}

	return b.String(), nil
}

func PrettyPrintDataFile(dataFileFullPath string, editText func(todo string, index int) string) error {
	result, err := SPrettyPrintDataFile(dataFileFullPath, editText)
	if err != nil {
		return fmt.Errorf("SPrettyPrintDataFile: %w", err)
	}
	fmt.Println(result)
	return nil
}
