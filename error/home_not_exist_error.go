package error

import "fmt"

type FilePathNotStartsWithHome struct {
	HomeStr string
}

func (e *FilePathNotStartsWithHome) Error() string {
	return fmt.Sprintf("File path does not start with %v", e.HomeStr)
}
