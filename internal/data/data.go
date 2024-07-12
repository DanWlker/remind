package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/DanWlker/remind/internal/config"
	"github.com/DanWlker/remind/internal/shared"
)

type EditTextFunc func(todo string, index int) (string, error)

// This does not create the file if it doesn't exist
func GetTodoFromFile(fileFullPath string) (items []TodoEntity, err error) {
	f, err := os.Open(fileFullPath)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	items, err = shared.FGetStructFromYaml[TodoEntity](f)
	if err != nil {
		return nil, fmt.Errorf("FGetStructFromYaml: %w", err)
	}
	return items, nil
}

func WriteTodoToFile(fileFullPath string, todoList []TodoEntity) (err error) {
	// Opens with
	// Write - WRONLY
	// Create if not exist - CREATE
	// Truncate (empty the file?) - Truncate
	f, err := os.OpenFile(fileFullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = errors.Join(
				err,
				fmt.Errorf("f.Close: %w", err2),
			)
		}
	}()

	if err = shared.FWriteStructToYaml(f, todoList); err != nil {
		return fmt.Errorf("FWriteStructToYaml: %w", err)
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

func SPrettyPrintFile(fileFullPath string, editText EditTextFunc) (string, error) {
	todoList, err := GetTodoFromFile(fileFullPath)
	if err != nil {
		return "", fmt.Errorf("GetTodoFromFile: %w", err)
	}

	var b bytes.Buffer
	if err := FPrettyPrintFile(&b, todoList, editText); err != nil {
		return "", fmt.Errorf("FPrettyPrintFile: %w", err)
	}
	return b.String(), nil
}

func PrettyPrintFile(fileFullPath string, editText EditTextFunc) error {
	todoList, err := GetTodoFromFile(fileFullPath)
	if err != nil {
		return fmt.Errorf("GetTodoFromFile: %w", err)
	}

	if err := FPrettyPrintFile(os.Stdout, todoList, editText); err != nil {
		return fmt.Errorf("FPrettyPrintFile: %w", err)
	}

	return nil
}

func FPrettyPrintFile(w io.Writer, todoList []TodoEntity, editText EditTextFunc) error {
	for i, todo := range todoList {
		if editText == nil {
			if _, err := fmt.Fprintln(w, todo.Text); err != nil {
				return fmt.Errorf("editText is nil, fmt.Fprintln: %w", err)
			}
			continue
		}
		res, err := editText(todo.Text, i)
		if err != nil {
			return fmt.Errorf("editText: %w", err)
		}
		if _, err := fmt.Fprintln(w, res); err != nil {
			return fmt.Errorf("editText not nil, fmt.Fprintln: %w", err)
		}
	}

	return nil
}
