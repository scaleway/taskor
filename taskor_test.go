package taskor

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/scaleway/taskor/mock"
)

func Test_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test fails at compilation time if there is a desynchro between mock and the library
	t.Run("Check that mocked taskor implements the taskor interface", func(t *testing.T) {
		mockTaskor := mock_taskor.NewMockTaskManager(ctrl)
		var _ TaskManager = mockTaskor
	})
}
