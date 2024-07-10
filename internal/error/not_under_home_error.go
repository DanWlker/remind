package error

import "fmt"

type NotUnderHomeError struct {
	HomeStr string
	FileStr string
}

func (e NotUnderHomeError) Error() string {
	return fmt.Sprintf("File path %v does not start with: %v", e.FileStr, e.HomeStr)
}
