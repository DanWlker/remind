package error

import "fmt"

type RecordDoesNotExistError struct {
	ID string
}

func (e RecordDoesNotExistError) Error() string {
	return fmt.Sprintf("Record requested does not exist, identifier: %v", e.ID)
}
