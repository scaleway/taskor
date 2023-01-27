package retry

import "time"

type countDownRetry struct {
	duration time.Duration
}

// CountDownRetry return an implementation of RetryMechanismFunc interface
func CountDownRetry(duration time.Duration) RetryMechanismFunc {
	return &countDownRetry{duration: duration}
}

// DurationBeforeRetry method to implement RetryMechanismFunc interface
// This return a duration defined during initialization of type
func (c countDownRetry) DurationBeforeRetry(currentTry int) time.Duration {
	// currentTry is not using for this retry mechanism
	return c.duration
}
