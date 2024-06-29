package shared

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatPathToRemoveHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir: %v", err)
	}
	tests := []struct {
		before, after         string
		errorStr              string
		checkValAfterErrMatch bool
	}{
		{home, ".", "", false},
		{filepath.Join(home, "celebrate"), "celebrate", "", false},
		{filepath.Join(home, "celebrate", "something"), filepath.Join("celebrate", "something"), "", false},
		{"", "", "does not start with", true},
		{"../../", "../../", "does not start with", true},
		{filepath.Join("celebrate", "something"), filepath.Join("celebrate", "something"), "does not start with", true},
	}

	for _, tt := range tests {
		t.Run(tt.before, func(t *testing.T) {
			result, err := FormatRemoveHome(tt.before)
			if err != nil {
				if !ErrorContains(err, tt.errorStr) {
					t.Errorf("expected error %v, got error %v", tt.errorStr, err)
					return
				}

				if !tt.checkValAfterErrMatch {
					return
				}
			}

			if result != tt.after {
				t.Errorf("expected %v, got %v", tt.after, result)
			}
		})
	}
}

func TestGetHomeRemovedWorkingDir(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}

	t.Run(dir, func(t *testing.T) {
		res, err := GetHomeRemovedWorkingDir()
		if err != nil {
			t.Errorf("Expected no errors, got %v", err)
		}

		if res == dir {
			t.Errorf("Home removed dir should not be the same as non home removed dir")
		}
	})
}
