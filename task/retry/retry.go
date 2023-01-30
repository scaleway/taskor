package retry

import "time"

// RetryMechanism interface to handling
// different way to waiting before retry a task
type RetryMechanism interface {
	DurationBeforeRetry(currentTry int) time.Duration
}
