package error

import "fmt"

type FilePathNotStartsWithHome struct {
}

func (e *FilePathNotStartsWithHome) Error() string {
	return fmt.Sprintf("File path does not start with $Home")
}
