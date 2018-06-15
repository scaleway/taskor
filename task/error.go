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
