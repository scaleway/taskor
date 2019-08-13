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
}

// New create a new Taskor instance
func New(runner runner.Runner) (TaskManager, error) {
	log.Debug("Starting")
	serializer.GlobalSerializer = serializer.TypeJSON
	taskor, err := handler.New(runner)
	return taskor, err
}

// SetLogger - change current logger
func SetLogger(newLogger log.Logger) {
	log.SetLogger(newLogger)
}
