package task

import (
	"time"

	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task/retry"
	"github.com/scaleway/taskor/utils"
)

const taskIDSize = 15

var (
	defaultRetryMechanism retry.RetryMechanismFunc = retry.CountDownRetry(20 * time.Second)
)

// Definition struct used to define task
type Definition struct {
	Name string
	Run  func(task *Task) error
}

func (d Definition) LoggerFields() map[string]interface{} {
	result := make(map[string]interface{})
	result["Name"] = d.Name
	return result
}

// Task struct used to be send in queue
type Task struct {
	// TaskID string (doesn't change on retry)
	ID string
	// RunningID Id of current running (change on retry)
	RunningID string
	// TaskName name of task to execute
	TaskName string
	// Parameter serialized task parameter
	Parameter []byte
	// Serialier Serializer to use to unserialize parameter
	Serializer serializer.Type
	// DateQueued date the task was queued
	DateQueued time.Time
	// DateExecuted date the task was executed
	DateExecuted time.Time
	// DateDone date the task was done (end of execution)
	DateDone time.Time
	// MaxRetry max retry allowed, negative value mean infinit
	MaxRetry int
	// CurrentTry (starts at 1)
	CurrentTry int
	// RetryOnError define is the task should retry if the task return err != nil
	RetryOnError bool
	// RetryMechanismFunc Interface to implement different method
	// to calculate duration to wait before retry
	RetryMechanismFunc retry.RetryMechanismFunc
	// ETA time after the task can be exec
	ETA time.Time
	// Error last error that was return by the task
	Error string
	// LinkError task
	LinkError *Task
	// ChildTasks Task
	ChildTasks []*Task
	// ParentTask access to the parent task
	ParentTask *Task
}

func (t Task) LoggerFields() map[string]interface{} {
	result := make(map[string]interface{})
	result["ID"] = t.ID
	result["RunningID"] = t.RunningID
	result["TaskName"] = t.TaskName
	result["MaxRetry"] = t.MaxRetry
	result["CurrentTry"] = t.CurrentTry

	if t.ParentTask != nil {
		result["ParentTask_ID"] = t.ParentTask.ID
		result["ParentTask_RunningID"] = t.ParentTask.RunningID
		result["ParentTask_Name"] = t.ParentTask.TaskName
	}

	if t.LinkError != nil {
		result["ErrorTask_Name"] = t.LinkError.TaskName
	}
	return result
}

// CreateTask create a new task without running it
func CreateTask(taskName string, param interface{}) (*Task, error) {
	// Serialize parameter
	serializedParameter, err := serializer.GetGlobalSerializer().Serialize(param)
	if err != nil {
		return nil, err
	}

	task := &Task{
		TaskName:   taskName,
		Parameter:  serializedParameter,
		Serializer: serializer.GlobalSerializer,
		CurrentTry: 0,
		// Default is don't retry
		MaxRetry: 0,
		// Wait 20 second before retry
		RetryMechanismFunc: defaultRetryMechanism,
		// Task can be exec starting now
		ETA: time.Now(),
		ID:  utils.GenerateRandString(taskIDSize),
	}
	return task, nil
}

// UnserializeParameter unserialize task parameter using task serializer
func (t *Task) UnserializeParameter(v interface{}) error {
	return serializer.GetSerializer(t.Serializer).Unserialize(v, t.Parameter)
}

// GetID return current task ID
func (t *Task) GetID() string {
	return t.ID
}

// SetMaxRetry Define a max retry
func (t *Task) SetMaxRetry(retry int) *Task {
	t.MaxRetry = retry
	return t
}

// SetCurrentTry return current try
func (t *Task) SetCurrentTry(v int) *Task {
	t.CurrentTry = v
	return t
}

// SetRetryOnError define retry strategie
func (t *Task) SetRetryOnError(v bool) *Task {
	t.RetryOnError = v
	return t
}

// SetCountDownRetry define time to wait before retry
func (t *Task) SetCountDownRetry(duration time.Duration) *Task {
	log.Warn("SetCountDownRetry function is deprecated: use SetRetryMechanism(...) instead")
	t.RetryMechanismFunc = retry.CountDownRetry(duration)
	return t
}

// SetRetryMechanism define algorithm to calculate duration to wait before retry
func (t *Task) SetRetryMechanism(mechanismFunc retry.RetryMechanismFunc) *Task {
	t.RetryMechanismFunc = mechanismFunc
	return t
}

// SetETA define time after that task can be exec
func (t *Task) SetETA(eta time.Time) *Task {
	t.ETA = eta
	return t
}

// SetLinkError define task that be call in error case
func (t *Task) SetLinkError(linkedErrorTask *Task) *Task {
	t.LinkError = linkedErrorTask
	return t
}

// AddChild Add a task that will be run after this one
func (t *Task) AddChild(childTask *Task) *Task {
	if childTask == nil {
		return t
	}
	t.ChildTasks = append(t.ChildTasks, childTask)
	return t
}

// LastRetry determines if no more retries are allowed
func (t *Task) LastRetry() bool {
	if t.MaxRetry == -1 {
		return false
	}
	return t.CurrentTry >= t.MaxRetry
}
