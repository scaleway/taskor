package retry

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	// ErrCountDownRetryInvalidParams Raise invalid params when a params is not found
	ErrCountDownRetryInvalidParams = fmt.Errorf("invalid params")
	// ErrCountDownRetryInvalidDuration Raise invalid duration when duration from params can't be parsed as Duration
	ErrCountDownRetryInvalidDuration = fmt.Errorf("invalid duration")
)

type countDownRetry struct {
	duration time.Duration
}

// Type return MechanismType
func (c *countDownRetry) Type() RetryMechanismType {
	return CountDownRetryMechanismType
}

// NewCountDownRetryFromDefinition initialize CountDownRetry from RetryMechanismDefinition
func NewCountDownRetryFromDefinition(definition RetryMechanismDefinition) (RetryMechanism, error) {
	durationStr, ok := definition.Params["duration"]
	if !ok {
		return nil, ErrCountDownRetryInvalidParams
	}

	duration, err := time.ParseDuration(durationStr.(string))
	if err != nil {
		return nil, ErrCountDownRetryInvalidDuration
	}

	return CountDownRetry(duration), nil
}

// MarshalJSON implement JSON Marshaller to encoding this complex object
func (c *countDownRetry) MarshalJSON() ([]byte, error) {
	return json.Marshal(RetryMechanismDefinition{
		Type: CountDownRetryMechanismType,
		Params: map[string]interface{}{
			"duration": c.duration.String(),
		},
	})
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
