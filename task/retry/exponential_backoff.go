package retry

import (
	"encoding/json"
	"fmt"
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

	// ErrExponentialBackOffRetryInvalidParams Raise invalid params when a params is not found
	ErrExponentialBackOffRetryInvalidParams = fmt.Errorf("invalid params")
	// ErrExponentialBackOffRetryInvalidDuration Raise invalid duration when duration from params can't be parsed as Duration
	ErrExponentialBackOffRetryInvalidDuration = fmt.Errorf("invalid duration")
)

type exponentialBackOffRetry struct {
	*backoff.Backoff
}

// Type return MechanismType
func (e *exponentialBackOffRetry) Type() RetryMechanismType {
	return ExponentialBackOffRetryMechanismType
}

// NewExponentialBackOffRetryFromDefinition initialize ExponentialBackOffRetry from RetryMechanismDefinition
func NewExponentialBackOffRetryFromDefinition(definition RetryMechanismDefinition) (RetryMechanism, error) {
	factorValue, ok := definition.Params["factor"]
	if !ok {
		return nil, ErrExponentialBackOffRetryInvalidParams
	}

	jitter, ok := definition.Params["jitter"]
	if !ok {
		return nil, ErrExponentialBackOffRetryInvalidParams
	}

	minDurationStr, ok := definition.Params["min_duration"]
	if !ok {
		return nil, ErrExponentialBackOffRetryInvalidParams
	}

	minDuration, err := time.ParseDuration(minDurationStr.(string))
	if err != nil {
		return nil, ErrExponentialBackOffRetryInvalidDuration
	}

	maxDurationStr, ok := definition.Params["max_duration"]
	if !ok {
		return nil, ErrExponentialBackOffRetryInvalidParams
	}

	maxDuration, err := time.ParseDuration(maxDurationStr.(string))
	if err != nil {
		return nil, ErrExponentialBackOffRetryInvalidDuration
	}

	var factor float64
	switch factorValue := factorValue.(type) {
	case int:
		factor = float64(factorValue)
	case float64:
		factor = factorValue
	default:
		return nil, ErrExponentialBackOffRetryInvalidParams
	}

	return ExponentialBackOffRetry(
		SetFactor(factor),
		SetJitter(jitter.(bool)),
		SetMin(minDuration),
		SetMax(maxDuration),
	), nil
}

// MarshalJSON implement JSON Marshalling to encode this complex object
func (e *exponentialBackOffRetry) MarshalJSON() ([]byte, error) {
	return json.Marshal(RetryMechanismDefinition{
		Type: ExponentialBackOffRetryMechanismType,
		Params: map[string]interface{}{
			"factor":       e.Factor,
			"jitter":       e.Jitter,
			"min_duration": e.Min.String(),
			"max_duration": e.Max.String(),
		},
	})
}

// ExponentialBackOffRetryOption Implement Option pattern
// This permit to define a list of option available to pass to the constructor
type ExponentialBackOffRetryOption func(retry *exponentialBackOffRetry)

// SetFactor define factor to use
// Factor is the multiplying factor for each increment step.
func SetFactor(factor float64) ExponentialBackOffRetryOption {
	return func(retry *exponentialBackOffRetry) {
		retry.Factor = factor
	}
}

// SetJitter define if jitter must be used
// Jitter eases contention by randomizing backoff steps.
func SetJitter(jitter bool) ExponentialBackOffRetryOption {
	return func(retry *exponentialBackOffRetry) {
		retry.Jitter = jitter
	}
}

// SetMin define minimum amount of time to wait before nex retry
func SetMin(min time.Duration) ExponentialBackOffRetryOption {
	return func(retry *exponentialBackOffRetry) {
		retry.Min = min
	}
}

// SetMax define maximum amount of time to wait before next retry
func SetMax(max time.Duration) ExponentialBackOffRetryOption {
	return func(retry *exponentialBackOffRetry) {
		retry.Max = max
	}
}

// ExponentialBackOffRetry return an implementation of RetryMechanism interface
func ExponentialBackOffRetry(opts ...ExponentialBackOffRetryOption) RetryMechanism {
	e := &exponentialBackOffRetry{defaultExponentialBackOffSettings}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// DurationBeforeRetry method to implement RetryMechanism interface
// This return a duration defined during initialization of type
func (e *exponentialBackOffRetry) DurationBeforeRetry(currentTry int) time.Duration {
	return e.ForAttempt(float64(currentTry - 1))
}
