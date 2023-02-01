package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewRetryMechanismFromDefinition(t *testing.T) {
	testCases := []struct {
		definition    RetryMechanismDefinition
		expected      RetryMechanism
		expectedError error
	}{
		{
			definition: RetryMechanismDefinition{
				Type: CountDownRetryMechanismType,
				Params: map[string]interface{}{
					"duration": "1h",
				},
			},
			expected: CountDownRetry(time.Hour * 1),
		},
		{
			definition: RetryMechanismDefinition{
				Type:   CountDownRetryMechanismType,
				Params: map[string]interface{}{},
			},
			expectedError: ErrCountDownRetryInvalidParams,
		},
		{
			definition: RetryMechanismDefinition{
				Type: CountDownRetryMechanismType,
				Params: map[string]interface{}{
					"duration": "invalid",
				},
			},
			expectedError: ErrCountDownRetryInvalidDuration,
		},
		{
			definition: RetryMechanismDefinition{
				Type: ExponentialBackOffRetryMechanismType,
				Params: map[string]interface{}{
					"factor":       2,
					"jitter":       true,
					"min_duration": "1m",
					"max_duration": "1h",
				},
			},
			expected: ExponentialBackOffRetry(
				SetFactor(2),
				SetJitter(true),
				SetMin(time.Minute*1),
				SetMax(time.Hour*1),
			),
		},
		{
			definition: RetryMechanismDefinition{
				Type: ExponentialBackOffRetryMechanismType,
				Params: map[string]interface{}{
					"min_duration": "1m",
					"max_duration": "1h",
				},
			},
			expected:      nil,
			expectedError: ErrExponentialBackOffRetryInvalidParams,
		},
		{
			definition: RetryMechanismDefinition{
				Type: ExponentialBackOffRetryMechanismType,
				Params: map[string]interface{}{
					"factor":       2,
					"jitter":       true,
					"min_duration": "invalid",
				},
			},
			expected:      nil,
			expectedError: ErrExponentialBackOffRetryInvalidDuration,
		},
		{
			definition: RetryMechanismDefinition{
				Type:   RetryMechanismType("UNKNOWN"),
				Params: nil,
			},
			expected:      nil,
			expectedError: ErrRetryMechanismTypeNotImplemented,
		},
	}

	for _, testCase := range testCases {
		rm, err := NewRetryMechanismFromDefinition(testCase.definition)
		assert.Equal(t, testCase.expected, rm)
		assert.Equal(t, testCase.expectedError, err)
	}

}
