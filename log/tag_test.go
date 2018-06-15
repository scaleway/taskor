package log

import (
	"github.com/scaleway/taskor/task"
	"reflect"
	"testing"
)

var fixtureTask, _ = task.CreateTask("test", nil)

func Test_GetFields(t *testing.T) {
	t.Run("simple task", func(t *testing.T) {
		expected := map[string]interface{}{
			"ID": fixtureTask.ID,
			"RunningID": fixtureTask.RunningID,
			"TaskName": fixtureTask.TaskName,
			"MaxRetry": fixtureTask.MaxRetry,
			"CurrentTry": fixtureTask.CurrentTry,
		}
		got := GetFields(*fixtureTask)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Taskor.GetFields() got %v, want %v", got, expected)
		}
	})

	t.Run("*Task with parents", func(t *testing.T) {
		childTask, _ := task.CreateTask("child", nil)
		childTask.ParentTask = fixtureTask
		expected := map[string]interface{}{
			"ID": childTask.ID,
			"RunningID": childTask.RunningID,
			"TaskName": childTask.TaskName,
			"MaxRetry": childTask.MaxRetry,
			"CurrentTry": childTask.CurrentTry,
			"ParentTask_ID": fixtureTask.ID,
			"ParentTask_RunningID": fixtureTask.RunningID,
			"ParentTask_Name": fixtureTask.TaskName,
		}
		got := GetFields(childTask)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Taskor.GetFields() got %v, want %v", got, expected)
		}
	})

	t.Run("Definition", func(t *testing.T) {
		definition := task.Definition{
			Name: "test",
			Run:  func(t *task.Task) error { return nil },
		}
		expected := map[string]interface{}{"Name": definition.Name}
		got := GetFields(definition)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Taskor.GetFields() got %v, want %v", got, expected)
		}
	})
}
