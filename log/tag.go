package log

import (
	"github.com/scaleway/taskor/task"
	"reflect"
)

// GetFields return a map[string]interface of exportable fields for log
func GetFields(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	sv := reflect.ValueOf(v)

	switch value := v.(type) {
	case task.Task:
		return getTaskFields(value)
	case task.Definition:
		return getDefinitionFields(value)
	}

	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		// If object is pointer and different nil recall function with value
		return GetFields(sv.Elem().Interface())
	}
	return result}

func getTaskFields(taskToLog task.Task) map[string]interface{} {
	result := make(map[string]interface{})
	result["ID"] = taskToLog.ID
	result["RunningID"] = taskToLog.RunningID
	result["TaskName"] = taskToLog.TaskName
	result["MaxRetry"] = taskToLog.MaxRetry
	result["CurrentTry"] = taskToLog.CurrentTry

	if taskToLog.ParentTask != nil {
		result["ParentTask_ID"] = taskToLog.ParentTask.ID
		result["ParentTask_RunningID"] = taskToLog.ParentTask.RunningID
		result["ParentTask_Name"] = taskToLog.ParentTask.TaskName
	}

	if taskToLog.LinkError != nil {
		result["ErrorTask_Name"] = taskToLog.LinkError.TaskName
	}
	return result
}

func getDefinitionFields(definition task.Definition) map[string]interface{} {
	result := make(map[string]interface{})
	result["Name"] = definition.Name
	return result
}