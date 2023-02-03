package task

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task/retry"
	"github.com/stretchr/testify/assert"
)

var fixtureTask, _ = CreateTask("test", nil)

func Test_Task_LoggerFields(t *testing.T) {
	t.Run("simple task", func(t *testing.T) {
		expected := map[string]interface{}{
			"ID":         fixtureTask.ID,
			"RunningID":  fixtureTask.RunningID,
			"TaskName":   fixtureTask.TaskName,
			"MaxRetry":   fixtureTask.MaxRetry,
			"CurrentTry": fixtureTask.CurrentTry,
		}
		got := (*fixtureTask).LoggerFields()
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Task.LoggerFields() got %v, want %v", got, expected)
		}
	})

	t.Run("*Task with parents", func(t *testing.T) {
		childTask, _ := CreateTask("child", nil)
		childTask.ParentTask = fixtureTask
		expected := map[string]interface{}{
			"ID":                   childTask.ID,
			"RunningID":            childTask.RunningID,
			"TaskName":             childTask.TaskName,
			"MaxRetry":             childTask.MaxRetry,
			"CurrentTry":           childTask.CurrentTry,
			"ParentTask_ID":        fixtureTask.ID,
			"ParentTask_RunningID": fixtureTask.RunningID,
			"ParentTask_Name":      fixtureTask.TaskName,
		}
		got := childTask.LoggerFields()
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Task.LoggerFields() got %v, want %v", got, expected)
		}
	})
}

func Test_Definition_LoggerFields(t *testing.T) {
	definition := Definition{
		Name: "test",
		Run:  func(t *Task) error { return nil },
	}
	expected := map[string]interface{}{"Name": definition.Name}
	got := definition.LoggerFields()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Definition.LoggerFields() got %v, want %v", got, expected)
	}
}

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

func Test_Task_Serialize(t *testing.T) {
	task, _ := CreateTask("t1", nil)

	t2, _ := CreateTask("t2", nil)
	t2.RetryMechanism = retry.ExponentialBackOffRetry(retry.SetJitter(false), retry.SetMin(time.Second*5), retry.SetMax(time.Minute*1), retry.SetFactor(1.5))
	task.AddChild(t2)

	t3, _ := CreateTask("t3", nil)
	t3.RetryMechanism = retry.ExponentialBackOffRetry(retry.SetJitter(false), retry.SetMin(time.Minute*5), retry.SetMax(time.Hour*1), retry.SetFactor(3))
	task.AddChild(t3)

	data, err := serializer.GetSerializer(task.Serializer).Serialize(task)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	log.Print(string(data))

	newTask := Task{}
	err = serializer.GetSerializer(task.Serializer).Unserialize(&newTask, data)
	assert.Nil(t, err)
	assert.Equal(t, newTask.ID, task.ID)
	assert.Len(t, newTask.ChildTasks, 2)
	assert.Equal(t, newTask.ChildTasks[0].ID, t2.ID)
	assert.Equal(t, newTask.ChildTasks[0].TaskName, t2.TaskName)
	assert.Equal(t, newTask.ChildTasks[0].MaxRetry, t2.MaxRetry)
	assert.Equal(t, newTask.ChildTasks[0].RetryMechanism, t2.RetryMechanism)
	assert.Equal(t, newTask.ChildTasks[0].RetryMechanism, t2.RetryMechanism)
	assert.Equal(t, newTask.ChildTasks[0].RetryMechanism.Type(), t2.RetryMechanism.Type())
	assert.Equal(t, newTask.ChildTasks[0].RetryMechanism.DurationBeforeRetry(5), t2.RetryMechanism.DurationBeforeRetry(5))
	assert.Equal(t, newTask.ChildTasks[1].ID, t3.ID)
	assert.Equal(t, newTask.ChildTasks[1].TaskName, t3.TaskName)
	assert.Equal(t, newTask.ChildTasks[1].MaxRetry, t3.MaxRetry)
	assert.Equal(t, newTask.ChildTasks[1].RetryMechanism, t3.RetryMechanism)
	assert.Equal(t, newTask.ChildTasks[1].RetryMechanism.Type(), t3.RetryMechanism.Type())
	assert.Equal(t, newTask.ChildTasks[1].RetryMechanism.DurationBeforeRetry(5), t3.RetryMechanism.DurationBeforeRetry(5))
}

func Test_UnserializeTask(t *testing.T) {
	data := []byte(`{"ID":"dnJGerKR2qE112Y","RunningID":"","TaskName":"t2","Parameter":"bnVsbA==","Serializer":0,"DateQueued":"0001-01-01T00:00:00Z","DateExecuted":"0001-01-01T00:00:00Z","DateDone":"0001-01-01T00:00:00Z","MaxRetry":10,"CurrentTry":0,"RetryOnError":true,"RetryMechanism":{"type":"ExponentialBackOffRetry","params":{"factor":3,"jitter":false,"max_duration":"2h0m0s","min_duration":"5m0s"}},"ETA":"2023-02-01T15:29:22.527005+01:00","Error":"","LinkError":null,"ChildTasks":null,"ParentTask":null}`)

	task := Task{}
	err := serializer.GetGlobalSerializer().Unserialize(&task, data)
	assert.Nil(t, err)

	assert.Equal(t, task.ID, "dnJGerKR2qE112Y")
	assert.Equal(t, task.TaskName, "t2")
	assert.Equal(t, task.MaxRetry, 10)
	assert.Equal(t, task.RetryOnError, true)
	assert.Equal(t, task.RetryMechanism.Type(), retry.ExponentialBackOffRetryMechanismType)
	assert.Equal(t, task.RetryMechanism.DurationBeforeRetry(5), time.Hour*2)
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

func Test_SetDefaultRetryMechanism(t *testing.T) {
	rm := retry.ExponentialBackOffRetry(retry.SetJitter(false), retry.SetMin(time.Minute*5), retry.SetMax(time.Hour*1), retry.SetFactor(3))
	SetDefaultRetryMechanism(rm)

	task, _ := CreateTask("test", nil)
	assert.Equal(t, rm, task.RetryMechanism)
}
