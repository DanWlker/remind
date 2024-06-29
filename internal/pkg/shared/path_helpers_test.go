package shared

import (
	"reflect"
	"testing"

	i_error "github.com/DanWlker/remind/internal/error"
)

func TestFormatPathToRemoveHome(t *testing.T) {
	home := "/something/hello"
	var tests = []struct {
		before, after string
		error         error
	}{
		{"", "", &i_error.FilePathNotStartsWithHome{}},
		{home + "", "", nil},
		{home + "/celebrate", "/celebrate", nil},
		{"/celebrate", "/celebrate", nil},
	}

	for _, tt := range tests {
		t.Run(tt.before, func(t *testing.T) {
			result, err := FormatPathToRemoveHome(tt.before, func() (string, error) { return home, nil })
			if tt.error != nil {
				if reflect.TypeOf(err) == reflect.TypeOf(tt.error) {
					return
				}
				t.Errorf("expected %v, got %v", tt.error, err)
			}

			if result != tt.after {
				t.Errorf("expected %v, got %v", tt.after, result)
			}
		})
	}
}
