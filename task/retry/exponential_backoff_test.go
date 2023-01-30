package retry

import (
	"testing"
	"time"

	"github.com/jpillora/backoff"
	"github.com/stretchr/testify/assert"
)

func Test_ExponentialBackOff_ExponentialBackOffRetry(t *testing.T) {
	testCases := []struct {
		options         []ExponentialBackOffRetryOptionFn
		expectedOptions *backoff.Backoff
	}{
		{
			options: nil,
			expectedOptions: &backoff.Backoff{
				Factor: 2,
				Jitter: false,
				Min:    100 * time.Millisecond,
				Max:    10 * time.Second,
			},
		},
		{
			options: []ExponentialBackOffRetryOptionFn{
				SetFactor(3),
				SetJitter(true),
			},
			expectedOptions: &backoff.Backoff{
				Factor: 3,
				Jitter: true,
				Min:    time.Millisecond * 100,
				Max:    time.Second * 10,
			},
		},
		{
			options: []ExponentialBackOffRetryOptionFn{
				SetFactor(4),
				SetJitter(false),
				SetMin(time.Second * 10),
				SetMax(time.Minute * 10),
			},
			expectedOptions: &backoff.Backoff{
				Factor: 4,
				Jitter: false,
				Min:    time.Second * 10,
				Max:    time.Minute * 10,
			},
		},
	}

	for _, testCase := range testCases {
		cdr := ExponentialBackOffRetry(testCase.options...)
		assert.IsType(t, cdr.(RetryMechanismFunc), cdr)
		assert.Equal(t, cdr.(*exponentialBackOffRetry).Backoff, testCase.expectedOptions)
	}
}

func Test_ExponentialBackOff_DurationBeforeRetry(t *testing.T) {
	var ebr interface{} = &exponentialBackOffRetry{
		Backoff: &backoff.Backoff{
			Factor: 2,
			Jitter: false,
			Min:    time.Second * 5,
			Max:    time.Minute * 1,
		},
	}
	assert.IsType(t, ebr.(RetryMechanismFunc), ebr)

	testCases := []struct {
		currentTry       int
		durationExpected time.Duration
	}{
		{currentTry: 0, durationExpected: time.Second * 5},
		{currentTry: 1, durationExpected: time.Second * 5},
		{currentTry: 2, durationExpected: time.Second * 10},
		{currentTry: 3, durationExpected: time.Second * 20},
		{currentTry: 4, durationExpected: time.Second * 40},
		{currentTry: 5, durationExpected: time.Minute * 1},
		{currentTry: 10, durationExpected: time.Minute * 1},
	}

	for _, testCase := range testCases {
		duration := ebr.(RetryMechanismFunc).DurationBeforeRetry(testCase.currentTry)
		assert.Equal(t, testCase.durationExpected, duration)
	}
}
