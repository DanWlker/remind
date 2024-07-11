package error

import "fmt"

type NotUnderHomeError struct {
	Home string
	File string
}

func (e NotUnderHomeError) Error() string {
	return fmt.Sprintf("file path %v does not start with: %v", e.File, e.Home)
}
