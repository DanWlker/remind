package data

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"

	"github.com/DanWlker/remind/internal/config"
)

// This does not create the file if it doesn't exist
func GetTodoFromFile(fileFullPath string) ([]TodoEntity, error) {
	file, err := os.ReadFile(fileFullPath)
	if err != nil {
		return []TodoEntity{}, fmt.Errorf("os.ReadFile: %w", err)
	}

	var items []TodoEntity
	if errUnmarshal := yaml.Unmarshal(file, &items); errUnmarshal != nil {
		return []TodoEntity{}, fmt.Errorf("yaml.Unmarshal: %w", errUnmarshal)
	}

	return items, nil
}

func WriteTodoToFile(fileFullPath string, todoList []TodoEntity) error {
	yamlTodoList, errMarshal := yaml.Marshal(todoList)
	if errMarshal != nil {
		return fmt.Errorf("yaml.Marshal: %w", errMarshal)
	}

	errWriteFile := os.WriteFile(fileFullPath, yamlTodoList, 0o644)
	if errWriteFile != nil {
		return fmt.Errorf("os.WriteFile: %w", errWriteFile)
	}
	return nil
}

func GetFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(config.UserDefinedDataFolder))
	if dataFolder == "" {
		home, errHomeDir := os.UserHomeDir()
		if errHomeDir != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", errHomeDir)
		}
		dataFolder = home + config.DefaultDataPathAfterHome
	}

	if errMkDirAll := os.MkdirAll(dataFolder, 0o770); errMkDirAll != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", errMkDirAll)
	}

	return dataFolder, nil
}

func SPrettyPrintFile(fileFullPath string, editText func(todo string, index int) string) (string, error) {
	var b bytes.Buffer
	todoList, errGetTodoFromDataFile := GetTodoFromFile(fileFullPath)
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

func PrettyPrintFile(fileFullPath string, editText func(todo string, index int) string) error {
	result, err := SPrettyPrintFile(fileFullPath, editText)
	if err != nil {
		return fmt.Errorf("SPrettyPrintDataFile: %w", err)
	}
	fmt.Println(result)
	return nil
}
