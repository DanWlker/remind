package error

import "fmt"

type RecordDoesNotExistError struct {
}

func (e *RecordDoesNotExistError) Error() string {
	return fmt.Sprintf("Record requested does not exist")
}
