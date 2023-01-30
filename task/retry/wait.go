package retry

import "time"

type countDownRetry struct {
	duration time.Duration
}

// CountDownRetry return an implementation of RetryMechanism interface
func CountDownRetry(duration time.Duration) RetryMechanism {
	return &countDownRetry{duration: duration}
}

// DurationBeforeRetry method to implement RetryMechanism interface
// This return a duration defined during initialization of type
func (c countDownRetry) DurationBeforeRetry(currentTry int) time.Duration {
	// currentTry is not using for this retry mechanism
	return c.duration
}
