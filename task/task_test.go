package task

import (
	"testing"
)

func Test_CreateTask(t *testing.T) {

	type args struct {
		taskName string
		param    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Task
		wantErr bool
	}{
		{
			name: "simple task",
			args: args{
				taskName: "task1",
				param:    nil,
			},
			want: Task{
				TaskName:   "task1",
				CurrentTry: 0,
				MaxRetry:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTask(tt.args.taskName, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Taskor.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.TaskName != tt.want.TaskName {
				t.Errorf("Taskor.CreateTask() = %v, want %v", got.TaskName, tt.want.TaskName)
			}
			if got.CurrentTry != tt.want.CurrentTry {
				t.Errorf("Taskor.CreateTask() = %v, want %v", got.CurrentTry, tt.want.CurrentTry)
			}
			if got.MaxRetry != tt.want.MaxRetry {
				t.Errorf("Taskor.CreateTask() = %v, want %v", got.MaxRetry, tt.want.MaxRetry)
			}
		})
	}
}

func TestTask_AddChild(t *testing.T) {
	t.Run("AddChild nil task", func(t *testing.T) {
		task, _ := CreateTask("test", nil)
		task.AddChild(nil)
		if len(task.ChildTasks) != 0 {
			t.Errorf("Task.AddChild() accept nil task")
		}
	})

}

func TestTask_LastRetry(t *testing.T) {
	tests := []struct {
		name   string
		define func(*Task) *Task
		want   bool
	}{
		{
			name: "no retries by default",
			define: func(task *Task) *Task {
				return task.SetMaxRetry(0)
			},
			want: true,
		},
		{
			name: "no retry limit",
			define: func(task *Task) *Task {
				return task.SetMaxRetry(-1)
			},
			want: false,
		},
		{
			name: "last retry",
			define: func(task *Task) *Task {
				return task.SetMaxRetry(2).SetCurrentTry(2)
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := CreateTask("test", nil)
			task = tt.define(task)
			if task.LastRetry() != tt.want {
				t.Errorf("Task.LastRetry() = %v, want %v", task.LastRetry(), tt.want)
			}
		})
	}
}
