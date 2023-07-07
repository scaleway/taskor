package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {
	message := "error message"
	metadata := map[string]string{
		"foo": "bar",
		"one": "two",
	}

	t.Run("WithoutOptions", func(t *testing.T) {
		expected := &Error{Message: "error message"}

		err := NewError(message)
		assert.Equal(t, expected, err)
	})

	t.Run("WithMetadata", func(t *testing.T) {
		expected := &Error{Message: "error message", Metadata: map[string]interface{}{
			"foo": "bar",
			"one": "two",
		}}

		err := NewError(message, WithMetadata(metadata))
		assert.Equal(t, expected, err)
	})
}
