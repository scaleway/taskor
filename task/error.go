package task

import (
	"errors"
)

var (
	// ErrTaskRetry error to return to retry a task
	ErrTaskRetry = errors.New("Retry task")
	// ErrNotRegisterd task has been pooled but was unknow (not register)
	ErrNotRegisterd = errors.New("Task was pooled but was not register")
)

// Error represent Task error
type Error struct {
	// Message error message
	Message string
	// Metadata context of error
	Metadata map[string]interface{}
}

// Error return Message of Error type
func (e *Error) Error() string {
	return e.Message
}

// ErrorOption implement option pattern
type ErrorOption func(taskErr *Error)

// WithMetadata enrich task error with metadata
func WithMetadata(metadata map[string]string) ErrorOption {
	return func(taskErr *Error) {
		md := make(map[string]interface{})
		for mdKey, mdValue := range metadata {
			md[mdKey] = mdValue
		}
		taskErr.Metadata = md
	}
}

// NewError initialize task.Error struct
func NewError(message string, errOpts ...ErrorOption) *Error {
	taskErr := &Error{Message: message}

	for _, errOpt := range errOpts {
		errOpt(taskErr)
	}

	return taskErr
}
