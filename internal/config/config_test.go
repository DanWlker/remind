package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigFolder(t *testing.T) {
	ex, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir: %v", err)
	}
	ex = filepath.Join(ex, "remind")

	t.Run(ex, func(t *testing.T) {
		res, err := GetConfigFolder()
		if err != nil {
			t.Errorf("Expected no errors, got %v", err)
		}

		if res != ex {
			t.Errorf("Expected %v, got %v", ex, res)
		}
	})
}
