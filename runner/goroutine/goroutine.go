package goroutine

import (
	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/task"
)

// RunnerConfig config use for goroutine runner
type RunnerConfig struct {
	MaxBufferedMessage int
	Concurrency        int
}

// Runner runner that use goroutine, can be used without separate worker
type Runner struct {
	internalChanTaskToRun chan task.Task
	config                RunnerConfig
}

// New Create a new Runner
func New(config RunnerConfig) *Runner {
	log.Warn("You are currently using a DEV/DEBUG runner please do not use it in production")
	g := Runner{}
	g.config = config
	return &g
}

// GetConcurrency retrieve concurrency settings for parallel task processing
func (g *Runner) GetConcurrency() int {
	return g.config.Concurrency
}

// Init channel
func (g *Runner) Init() error {
	g.internalChanTaskToRun = make(chan task.Task, g.config.MaxBufferedMessage)
	return nil
}

// Stop will be call when StopWorker will be call
func (g *Runner) Stop() error {
	close(g.internalChanTaskToRun)
	return nil
}

// Send send a new task to the pool
func (g *Runner) Send(t *task.Task) error {
	g.internalChanTaskToRun <- *t
	return nil
}

// RunWorkerTaskProvider runner that consume queue and push task to taskToRun chan
func (g *Runner) RunWorkerTaskProvider(taskToRun chan task.Task, stop <-chan bool) error {
loop:
	for {
		select {
		case <-stop:
			break loop
		case t := <-g.internalChanTaskToRun:
			taskToRun <- t
		}
	}
	return nil
}

// RunWorkerTaskAck runner that ack message when a task is done. Should stop on chan close
func (g *Runner) RunWorkerTaskAck(taskDone <-chan task.Task) {
	for {

		_, ok := <-taskDone
		if !ok {
			break
		}
	}
	return
}
