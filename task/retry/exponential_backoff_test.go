package retry

import (
	"testing"
	"time"

	"github.com/jpillora/backoff"
	"github.com/stretchr/testify/assert"
)

func Test_ExponentialBackOff_ExponentialBackOffRetry(t *testing.T) {
	testCases := []struct {
		options         []ExponentialBackOffRetryOption
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
			options: []ExponentialBackOffRetryOption{
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
			options: []ExponentialBackOffRetryOption{
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
		assert.IsType(t, cdr.(RetryMechanism), cdr)
		assert.Equal(t, cdr.(*exponentialBackOffRetry).Backoff, testCase.expectedOptions)
	}
}

func Test_ExponentialBackOff_Type(t *testing.T) {
	rm := ExponentialBackOffRetry()
	assert.Equal(t, ExponentialBackOffRetryMechanismType, rm.Type())
}

func Test_ExponentialBackOff_NewExponentialBackOffRetryFromDefinition(t *testing.T) {
	definition := RetryMechanismDefinition{
		Type: ExponentialBackOffRetryMechanismType,
		Params: map[string]interface{}{
			"factor":       1.1,
			"jitter":       false,
			"min_duration": "10m",
			"max_duration": "1h",
		},
	}

	rm, err := NewExponentialBackOffRetryFromDefinition(definition)
	assert.Nil(t, err)
	assert.Equal(t, rm.(*exponentialBackOffRetry).Factor, 1.1)
	assert.Equal(t, rm.(*exponentialBackOffRetry).Jitter, false)
	assert.Equal(t, rm.(*exponentialBackOffRetry).Min, 10*time.Minute)
	assert.Equal(t, rm.(*exponentialBackOffRetry).Max, 1*time.Hour)
}

func Test_ExponentialBackOff_MarshalJSON(t *testing.T) {
	rm := ExponentialBackOffRetry(
		SetFactor(3),
		SetJitter(true),
		SetMin(time.Millisecond*250),
		SetMax(time.Hour*2),
	)
	data, err := rm.MarshalJSON()

	expected := `{"type":"ExponentialBackOffRetry","params":{"factor":3,"jitter":true,"max_duration":"2h0m0s","min_duration":"250ms"}}`

	assert.Nil(t, err)
	assert.Equal(t, string(data), expected)
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
	assert.IsType(t, ebr.(RetryMechanism), ebr)

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
		duration := ebr.(RetryMechanism).DurationBeforeRetry(testCase.currentTry)
		assert.Equal(t, testCase.durationExpected, duration)
	}
}
