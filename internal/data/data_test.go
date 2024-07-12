package data

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestFGetTodoFromReader(t *testing.T) {
	t.Run("Test EOF returns without error", func(t *testing.T) {
		reader := strings.NewReader("")

		res, err := FGetTodoFromReader(reader)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(res) != 0 {
			t.Errorf("Expected empty array, got %v", res)
		}
	})
}

func TestFWriteTodoToFile(t *testing.T) {
	t.Run("Test nil array encodes without error", func(t *testing.T) {
		ex := "[]\n"

		var b bytes.Buffer

		err := FWriteTodoToFile(&b, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got := b.String()
		if got != ex {
			t.Errorf("Expected %v, got %v", strconv.Quote(ex), strconv.Quote(got))
		}
	})
}

func TestFPrettyPrintFile(t *testing.T) {
	t.Run("Functionality test", func(t *testing.T) {
		todoList := []TodoEntity{
			{"Text one"},
			{"Text two"},
		}
		ex := "Text one\nText two\n"

		var b bytes.Buffer
		err := FPrettyPrintFile(&b, todoList, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got := b.String()
		if got != ex {
			t.Errorf("Expected %v, got %v", strconv.Quote(ex), strconv.Quote(got))
		}
	})

	t.Run("Test nil array and nil function prints without error", func(t *testing.T) {
		var b bytes.Buffer
		err := FPrettyPrintFile(&b, nil, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got := b.String()
		if got != "" {
			t.Errorf("Expected empty string, got %v", strconv.Quote(got))
		}
	})

	t.Run("Test editText", func(t *testing.T) {
		todoList := []TodoEntity{
			{"Text one"},
			{"Text two"},
		}
		ex := "0: Text one\n1: Text two\n"

		var b bytes.Buffer
		err := FPrettyPrintFile(&b, todoList, func(todo string, index int) (string, error) {
			return fmt.Sprintf("%v: %v", index, todo), nil
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		got := b.String()
		if got != ex {
			t.Errorf("Expected %v, got %v", strconv.Quote(ex), strconv.Quote(got))
		}
	})
}
