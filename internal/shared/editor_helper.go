package shared

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var editor = "vim"

func init() {
	if s := os.Getenv("VISUAL"); s != "" {
		editor = s
		return
	}

	if s := os.Getenv("EDITOR"); s != "" {
		editor = s
		return
	}
}

func OpenDefaultEditor(data []byte) ([]byte, error) {
	return openEditor(editor, data)
}

func openEditor(editor string, data []byte) (result []byte, err error) {
	tempFile, err := os.CreateTemp("", "redit")
	if err != nil {
		return nil, fmt.Errorf("os.CreateTemp: %w", err)
	}
	defer func() {
		if err2 := os.Remove(tempFile.Name()); err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	// .Write will be a no-op if data is empty
	if _, err := tempFile.Write(data); err != nil {
		return nil, fmt.Errorf("tempFile.Write: %w", err)
	}

	cmd := exec.Command("sh", "-c", editor+" "+tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cmd.Run: %w", err)
	}

	result, err = os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return result, nil
}
