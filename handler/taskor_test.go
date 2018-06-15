package handler

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	runnerMock "github.com/scaleway/taskor/runner/mock"
	"github.com/scaleway/taskor/task"
)

var taskTest = task.Definition{
	Name: "test",
	Run:  func(t *task.Task) error { return nil },
}

var taskOtherTest = task.Definition{
	Name: "testOther",
	Run:  func(t *task.Task) error { return nil },
}

func TestTaskor_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()

	tests := []struct {
		name    string
		tasks   []*task.Definition
		wantErr bool
	}{
		{
			name:  "single task",
			tasks: []*task.Definition{&taskTest},
		},
		{
			name:    "double same task",
			tasks:   []*task.Definition{&taskTest, &taskTest},
			wantErr: true,
		},
		{
			name:    "double different task",
			tasks:   []*task.Definition{&taskTest, &taskOtherTest},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			taskManager, _ := New(mockRunner)
			for _, task := range tt.tasks {
				err = taskManager.Handle(task)
			}
			// Only check error on last one
			if err != nil != tt.wantErr {
				t.Errorf("Taskor.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestTaskor_GetHandled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()

	tests := []struct {
		name  string
		tasks []*task.Definition
		want  []*task.Definition
	}{
		{
			name:  "single task",
			tasks: []*task.Definition{&taskTest},
			want:  []*task.Definition{&taskTest},
		},
		{
			name:  "double different task",
			tasks: []*task.Definition{&taskTest, &taskOtherTest},
			want:  []*task.Definition{&taskTest, &taskOtherTest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			taskManager, _ := New(mockRunner)
			for _, task := range tt.tasks {
				taskManager.Handle(task)
			}

			got := taskManager.GetHandled()
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskManager.GetHandled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskor_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().Send(gomock.Any())
	taskManager, _ := New(mockRunner)

	t.Run("check task param", func(t *testing.T) {
		testTask, _ := task.CreateTask("test", nil)
		taskManager.Send(testTask)
		if testTask.ID == "" {
			t.Errorf("Task ID is nil")
		}

		if testTask.RunningID == "" {
			t.Errorf("Task RunningID is nil")
		}

		if time.Time.IsZero(testTask.DateQueued) {
			t.Errorf("Task DateQueued is nil")
		}
	})

}
