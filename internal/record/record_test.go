package record

import (
	"errors"
	"testing"

	i_error "github.com/DanWlker/remind/internal/error"
)

func TestFGetRecordIdentityWithIdentifier(t *testing.T) {
	t.Run("Functionality test", func(t *testing.T) {
		ex := RecordEntity{"up", "something"}
		items := []RecordEntity{
			{"never", "gonna"},
			{"give", "you"},
			ex,
		}

		got, err := FGetRecordIdentityWithIdentifier(items, "something")
		if err != nil {
			t.Errorf("Expectad no errors, got %v", err)
		}

		if got != ex {
			t.Errorf("Wrong record found, expected %v, got %v", ex, got)
		}
	})

	t.Run("Test record does not exist", func(t *testing.T) {
		var items []RecordEntity
		var ex i_error.RecordDoesNotExistError

		_, err := FGetRecordIdentityWithIdentifier(items, "something")
		if !errors.As(err, &ex) {
			t.Errorf("Should return %v, got %v", ex, err)
		}
	})
}
