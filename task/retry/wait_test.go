package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CountDownRetry_CountDownRetry(t *testing.T) {
	cdr := CountDownRetry(time.Minute * 15)
	assert.IsType(t, cdr.(RetryMechanism), cdr)

	duration := cdr.(RetryMechanism).DurationBeforeRetry(0)
	assert.Equal(t, time.Minute*15, duration)
}

func Test_CountDownRetry_DurationBeforeRetry(t *testing.T) {
	cdr := &countDownRetry{duration: time.Minute * 10}

	for _, currentRetry := range []int{1, 3, 5} {
		duration := cdr.DurationBeforeRetry(currentRetry)
		assert.Equal(t, time.Minute*10, duration)
	}
}

func Test_CountDownRetry_Type(t *testing.T) {
	rm := CountDownRetry(5 * time.Minute)
	assert.Equal(t, CountDownRetryMechanismType, rm.Type())
}

func Test_CountDownRetry_NewCountDownRetryFromDefinition(t *testing.T) {
	definition := RetryMechanismDefinition{
		Type: CountDownRetryMechanismType,
		Params: map[string]interface{}{
			"duration": "5m",
		},
	}

	rm, err := NewCountDownRetryFromDefinition(definition)
	assert.Nil(t, err)
	assert.Equal(t, rm.(*countDownRetry).duration, 5*time.Minute)
}

func Test_CountDownRetry_MarshalJSON(t *testing.T) {
	rm := CountDownRetry(time.Second * 10)
	data, err := rm.MarshalJSON()

	expected := `{"type":"CountDownRetry","params":{"duration":"10s"}}`

	assert.Nil(t, err)
	assert.Equal(t, string(data), expected)
}
