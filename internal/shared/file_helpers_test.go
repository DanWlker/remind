package shared

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

type SomeYamlStruct struct {
	Text string
}

func TestFGetStructFromReader(t *testing.T) {
	t.Run("Test EOF returns without error", func(t *testing.T) {
		reader := strings.NewReader("")

		res, err := FGetStructFromYaml[SomeYamlStruct](reader)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(res) != 0 {
			t.Errorf("Expected empty array, got %v", res)
		}
	})
}

func TestFWriteStructToFile(t *testing.T) {
	t.Run("Test nil array encodes without error", func(t *testing.T) {
		ex := "[]\n"

		var b bytes.Buffer

		err := FWriteStructToYaml[SomeYamlStruct](&b, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got := b.String()
		if got != ex {
			t.Errorf("Expected %v, got %v", strconv.Quote(ex), strconv.Quote(got))
		}
	})
}
