package error

import "fmt"

type RecordDoesNotExistError struct {
	RecordIdentifier string
}

func (e *RecordDoesNotExistError) Error() string {
	return fmt.Sprintf("Record requested does not exist, identifier: %v", e.RecordIdentifier)
}
