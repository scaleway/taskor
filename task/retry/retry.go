package retry

import "time"

// RetryMechanismFunc interface to support
// different delays before retrying a task
type RetryMechanismFunc interface {
	DurationBeforeRetry(currentTry int) time.Duration
}
