package shared

import (
	"fmt"
	"os"
	"os/exec"
)

var editor = "vim"

func init() {
	if s := os.Getenv("VISUAL"); s != "" {
		editor = s
	} else if s := os.Getenv("EDITOR"); s != "" {
		editor = s
	}
}

func OpenDefaultEditor(data []byte) ([]byte, error) {
	return openEditor(editor, data)
}

func openEditor(editor string, data []byte) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "redit")
	if err != nil {
		return nil, fmt.Errorf("os.CreateTemp: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		return nil, fmt.Errorf("tempFile.WriteString: %w", err)
	}

	cmd := exec.Command("sh", "-c", editor+" "+tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cmd.Run: %w", err)
	}

	result, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return result, nil
}
