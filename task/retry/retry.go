package retry

import (
	"fmt"
	"time"
)

var (
	// CountDownRetryMechanismType ...
	CountDownRetryMechanismType RetryMechanismType = "CountDownRetry"
	// ExponentialBackOffRetryMechanismType ...
	ExponentialBackOffRetryMechanismType RetryMechanismType = "ExponentialBackOffRetry"

	// ErrRetryMechanismTypeNotImplemented Raise error when mechanism type is not found in mechanismTypes
	ErrRetryMechanismTypeNotImplemented = fmt.Errorf("this retry mechanism type is not implemented")
)

var mechanismTypes = map[RetryMechanismType]func(definition RetryMechanismDefinition) (RetryMechanism, error){
	CountDownRetryMechanismType:          NewCountDownRetryFromDefinition,
	ExponentialBackOffRetryMechanismType: NewExponentialBackOffRetryFromDefinition,
}

// RetryMechanism interface to handling
// different way to waiting before retry a task
type RetryMechanism interface {
	Type() RetryMechanismType
	DurationBeforeRetry(currentTry int) time.Duration
	MarshalJSON() (b []byte, err error)
}

// RetryMechanismType ...
type RetryMechanismType string

// RetryMechanismDefinition ...
type RetryMechanismDefinition struct {
	Type   RetryMechanismType     `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// NewRetryMechanismFromDefinition initialize RetryMechanism interface
// from a given definition
func NewRetryMechanismFromDefinition(definition RetryMechanismDefinition) (RetryMechanism, error) {
	constructor, ok := mechanismTypes[definition.Type]
	if !ok {
		return nil, ErrRetryMechanismTypeNotImplemented
	}

	return constructor(definition)
}
