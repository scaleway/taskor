package retry

import "time"

// RetryMechanismFunc interface to handling
// different way to waiting before retry a task
type RetryMechanismFunc interface {
	DurationBeforeRetry(currentTry int) time.Duration
}
