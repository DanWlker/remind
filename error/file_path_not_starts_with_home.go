package error

import "fmt"

type FilePathNotStartsWithHome struct {
	HomeStr string
	FileStr string
}

func (e *FilePathNotStartsWithHome) Error() string {
	return fmt.Sprintf("File path %v does not start with: %v", e.FileStr, e.HomeStr)
}
