package shared

import (
	"fmt"
	"os"
	"os/exec"
)

var editor = "vim"

func init() {
	if s := os.Getenv("EDITOR"); s != "" {
		editor = s
	}
}

func OpenDefaultEditor(data []byte) ([]byte, error) {
	return openEditor(editor, data)
}

func openEditor(editor string, data []byte) ([]byte, error) {
	tempFile, errOsCreateTemp := os.CreateTemp("", "redit")
	if errOsCreateTemp != nil {
		return nil, fmt.Errorf("os.CreateTemp: %w", errOsCreateTemp)
	}
	defer os.Remove(tempFile.Name())

	if len(data) != 0 {
		_, errWriteString := tempFile.Write(data)
		if errWriteString != nil {
			return nil, fmt.Errorf("tempFile.WriteString: %w", errWriteString)
		}
	}

	cmd := exec.Command("sh", "-c", editor+" "+tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	errCmdRun := cmd.Run()
	if errCmdRun != nil {
		return nil, fmt.Errorf("cmd.Run: %w", errCmdRun)
	}

	result, errOsReadFile := os.ReadFile(tempFile.Name())
	if errOsReadFile != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", errOsReadFile)
	}

	return result, nil
}
