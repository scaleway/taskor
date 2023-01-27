package retry

import (
	"time"

	"github.com/jpillora/backoff"
)

var (
	defaultExponentialBackOffSettings = &backoff.Backoff{
		Factor: 2,
		Jitter: false,
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
	}
)

type exponentialBackOffRetry struct {
	*backoff.Backoff
}

// ExponentialBackOffRetryOptionFn Function to implement Option pattern
// This permit to define a list of option available to pass to the constructor
type ExponentialBackOffRetryOptionFn func(retry *exponentialBackOffRetry)

// SetFactor define factor to use
// Factor is the multiplying factor for each increment step.
func SetFactor(factor float64) ExponentialBackOffRetryOptionFn {
	return func(retry *exponentialBackOffRetry) {
		retry.Factor = factor
	}
}

// SetJitter define if jitter must be used
// Jitter eases contention by randomizing backoff steps.
func SetJitter(jitter bool) ExponentialBackOffRetryOptionFn {
	return func(retry *exponentialBackOffRetry) {
		retry.Jitter = jitter
	}
}

// SetMin define minimum amount of time to wait before nex retry
func SetMin(min time.Duration) ExponentialBackOffRetryOptionFn {
	return func(retry *exponentialBackOffRetry) {
		retry.Min = min
	}
}

// SetMax define maximum amount of time to wait before next retry
func SetMax(max time.Duration) ExponentialBackOffRetryOptionFn {
	return func(retry *exponentialBackOffRetry) {
		retry.Max = max
	}
}

// ExponentialBackOffRetry return an implementation of RetryMechanismFunc interface
func ExponentialBackOffRetry(opts ...ExponentialBackOffRetryOptionFn) RetryMechanismFunc {
	e := &exponentialBackOffRetry{defaultExponentialBackOffSettings}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// DurationBeforeRetry method to implement RetryMechanismFunc interface
// This return a duration defined during initialization of type
func (e *exponentialBackOffRetry) DurationBeforeRetry(currentTry int) time.Duration {
	return e.ForAttempt(float64(currentTry - 1))
}
