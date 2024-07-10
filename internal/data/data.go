package data

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/DanWlker/remind/internal/config"
)

// This does not create the file if it doesn't exist
func GetTodoFromFile(filename string) ([]TodoEntity, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	var (
		items []TodoEntity
		dec   = yaml.NewDecoder(f)
	)

	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return items, nil
}

func WriteTodoToFile(filename string, todoList []TodoEntity) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	if err := enc.Encode(todoList); err != nil {
		return fmt.Errorf("yaml.Marshal: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("closing file: %w", err)
	}

	return nil
}

func GetFolder() (string, error) {
	dataFolder := strings.TrimSpace(viper.GetString(config.UserDefinedDataFolder))
	if dataFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", err)
		}
		dataFolder = filepath.Join(home, config.DefaultDataSubdir)
	}

	if err := os.MkdirAll(dataFolder, 0o770); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	return dataFolder, nil
}

func SPrettyPrintFile(fileFullPath string, editText func(todo string, index int) string) (string, error) {
	var b bytes.Buffer

	if err := FPrettyPrintFile(&b, fileFullPath, editText); err != nil {
		return "", fmt.Errorf("FPrettyPrintDataFile: %w", err)
	}

	return b.String(), nil
}

func PrettyPrintFile(fileFullPath string, editText func(todo string, index int) string) error {
	return FPrettyPrintFile(os.Stdout, fileFullPath, editText)
}

func FPrettyPrintFile(w io.Writer, filename string, editText func(string, int) string) error {
	todoList, err := GetTodoFromFile(filename)
	if err != nil {
		return fmt.Errorf("GetTodoFromDataFile: %w", err)
	}

	for i, todo := range todoList {
		if editText == nil {
			fmt.Fprintln(w, todo.Text)
		} else {
			fmt.Fprintln(w, editText(todo.Text, i))
		}
	}

	return nil
}
