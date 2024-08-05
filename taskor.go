package taskor

import (
	"github.com/scaleway/taskor/handler"
	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/runner"
	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task"
)

// TaskManager Interface to communicate with client
type TaskManager interface {
	// Send a new task in queue
	Send(task *task.Task) error
	// Add a new task definition to be handle by worker
	Handle(Definition *task.Definition) error
	// Get all task definition that be handle
	GetHandled() []*task.Definition
	// Start to execute task in queue
	RunWorker() error
	// Stop worker
	StopWorker()
	// GetMetrics return current metric
	GetMetrics() handler.Metric
	// IsRunnerHealthy checks that the runner connection and channel are set
	IsRunnerHealthy() error
}

// New create a new Taskor instance
func New(runner runner.Runner) (TaskManager, error) {
	return NewWithSerializer(runner, serializer.TypeJSON)
}

func NewWithSerializer(runner runner.Runner, serializerType serializer.Type) (TaskManager, error) {
	log.Debug("Starting")
	return handler.NewWithSerializer(runner, serializerType)
}

// SetLogger - change current logger
func SetLogger(newLogger log.Logger) {
	log.SetLogger(newLogger)
}
