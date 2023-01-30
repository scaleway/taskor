package handler

import (
	"errors"
	"sync"
	"time"

	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/runner"
	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task"
	"github.com/scaleway/taskor/utils"
)

// Taskor implementation of TaskManager
type Taskor struct {
	runner   runner.Runner
	taskList map[string]*task.Definition

	// taskToRun is the chan used when a task need to be run.
	// Task will be analyze to know when it can be process
	taskToRun chan task.Task
	// taskToProcess is the chan to process the task directly without
	// any check (ETA check)
	taskToProcess chan task.Task
	// taskToSend is the chan used to send tasks to the queue
	taskToSend chan task.Task
	// taskDone is the chan used to inform task is done and can be ack
	taskDone chan task.Task

	// WaitGroup Use to know when goroutine are stop
	runWorkerTaskProviderWG sync.WaitGroup
	runWorkerTaskAckWG      sync.WaitGroup
	handlerTaskToRunWG      sync.WaitGroup
	handlerTaskToProcessWG  sync.WaitGroup
	handlerTaskToSendWG     sync.WaitGroup

	// stopWorkerTaskProvider chan use to stop taskprovider routine
	stopWorkerTaskProvider   chan bool
	stopHandlerTaskToRun     chan bool
	stopHandlerTaskToProcess chan bool
	stopHandlerTaskToSend    chan bool

	// boolean used to avoid stop a stopped worker
	workerRunning   bool
	workerStopMutex sync.Mutex

	// Metric
	metric Metric
}

// New create a new Taskor instance
func New(runner runner.Runner) (*Taskor, error) {
	// Init serializer
	serializer.GlobalSerializer = serializer.TypeJSON

	var t Taskor
	// Init task runner
	t.runner = runner
	err := t.runner.Init()
	if err != nil {
		return nil, err
	}
	// Init task list
	t.taskList = make(map[string]*task.Definition)
	t.metric = Metric{}
	return &t, nil
}

// Send send quickly a new task to the pool
func (t *Taskor) Send(taskToSend *task.Task) error {
	taskToSend.RunningID = utils.GenerateRandString(utils.TaskRunningIDSize)
	// Update queued date
	taskToSend.DateQueued = time.Now()
	log.InfoWithFields("Send task", taskToSend.LoggerFields())
	t.metric.TaskSent++
	return t.runner.Send(taskToSend)
}

// Handle register task that can be run
func (t *Taskor) Handle(definition *task.Definition) error {
	if _, ok := t.taskList[definition.Name]; ok {
		log.ErrorWithFields("Task name was already register", definition.LoggerFields())
		return errors.New("Task name was already register")
	}
	t.taskList[definition.Name] = definition
	return nil
}

// GetHandled return list of all task name registered
func (t *Taskor) GetHandled() []*task.Definition {
	handled := make([]*task.Definition, 0, len(t.taskList))
	for _, def := range t.taskList {
		handled = append(handled, def)
	}
	return handled
}

// GetMetrics return a copy of actual metrics
func (t *Taskor) GetMetrics() Metric {
	return t.metric
}
