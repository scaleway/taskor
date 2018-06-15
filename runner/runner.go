package runner

import (
	"github.com/scaleway/taskor/task"
)

// Runner Task runner
type Runner interface {
	// Init This method is call when TaskManager is created
	Init() error
	// Stop will be call when StopWorker will be call
	Stop() error
	// Send send a new task to the pool
	Send(*task.Task) error
	// RunWorkerTaskProvider runner that consume queue and push task to taskToRun chan
	RunWorkerTaskProvider(taskToRun chan task.Task, stop <-chan bool) error
	// RunWorkerTaskAck runner that ack message when a task is done. Should stop on chan close
	RunWorkerTaskAck(taskDone <-chan task.Task)
}
