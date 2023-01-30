package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CountDownRetry_CountDownRetry(t *testing.T) {
	cdr := CountDownRetry(time.Minute * 15)
	assert.IsType(t, cdr.(RetryMechanismFunc), cdr)

	duration := cdr.(RetryMechanismFunc).DurationBeforeRetry(0)
	assert.Equal(t, time.Minute*15, duration)
}

func Test_CountDownRetry_DurationBeforeRetry(t *testing.T) {
	cdr := &countDownRetry{duration: time.Minute * 10}

	for _, currentRetry := range []int{1, 3, 5} {
		duration := cdr.DurationBeforeRetry(currentRetry)
		assert.Equal(t, time.Minute*10, duration)
	}

}
