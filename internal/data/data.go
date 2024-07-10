package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"

	"github.com/DanWlker/remind/internal/config"
)

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

	dec := yaml.NewDecoder(f)
	err = dec.Decode(&items)
	if errors.Is(err, io.EOF) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("dec.Decode: %w", err)
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

	enc := yaml.NewEncoder(f)
	if err := enc.Encode(todoList); err != nil {
		return fmt.Errorf("enc.Encode: %w", err)
	}

	return err
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
		return "", fmt.Errorf("FPrettyPrintFile: %w", err)
	}

	return b.String(), nil
}

func PrettyPrintFile(fileFullPath string, editText func(todo string, index int) string) error {
	if err := FPrettyPrintFile(os.Stdout, fileFullPath, editText); err != nil {
		return fmt.Errorf("FPrettyPrintFile: %w", err)
	}
	return nil
}

func FPrettyPrintFile(w io.Writer, fileName string, editText func(todo string, index int) string) error {
	todoList, err := GetTodoFromFile(fileName)
	if err != nil {
		return fmt.Errorf("GetTodoFromFile: %w", err)
	}

	for i, todo := range todoList {
		if editText == nil {
			if _, err := fmt.Fprintln(w, todo.Text); err != nil {
				return fmt.Errorf("editText is nil, fmt.Fprintln: %w", err)
			}
			continue
		}
		if _, err := fmt.Fprintln(w, editText(todo.Text, i)); err != nil {
			return fmt.Errorf("editText not nil, fmt.Fprintln: %w", err)
		}
	}

	return nil
}
