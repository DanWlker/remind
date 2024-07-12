package shared

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

func FGetStructFromYaml[T any](r io.Reader) ([]T, error) {
	var items []T
	dec := yaml.NewDecoder(r)

	err := dec.Decode(&items)
	if errors.Is(err, io.EOF) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("dec.Decode: %w", err)
	}

	return items, nil
}

func FWriteStructToYaml[T any](w io.Writer, todoList []T) error {
	enc := yaml.NewEncoder(w)
	if err := enc.Encode(todoList); err != nil {
		return fmt.Errorf("enc.Encode: %w", err)
	}

	return nil
}
