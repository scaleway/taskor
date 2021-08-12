package handler

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/task"
)

var errorWorkerAlreadyRunning = errors.New("worker is already start")

// RunWorker run worker that wait new task and exec
func (t *Taskor) RunWorker() error {
	if t.workerRunning {
		return errorWorkerAlreadyRunning
	}
	t.workerRunning = true

	// taskToRun is the chan used when a task need to be run.
	// Task will be analyze to know when it can be process
	t.taskToRun = make(chan task.Task)
	// taskToProcess is the chan to process the task directly without
	// any check (ETA check)
	t.taskToProcess = make(chan task.Task)
	// taskToSend is the chan used to send tasks to the queue
	t.taskToSend = make(chan task.Task)
	// taskDone is the chan used to inform task is done and can be ack
	t.taskDone = make(chan task.Task)
	// stopChan Are chan use to stop all goroutine
	t.stopWorkerTaskProvider = make(chan bool)
	t.stopHandlerTaskToProcess = make(chan bool)
	t.stopHandlerTaskToSend = make(chan bool)
	t.stopHandlerTaskToRun = make(chan bool)

	// Handling SIGTERM & SIGINT
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		t.StopWorker()
	}()

	t.runWorkerTaskProviderWG.Add(1)
	go func() {
		t.runner.RunWorkerTaskProvider(t.taskToRun, t.stopWorkerTaskProvider)
		t.runWorkerTaskProviderWG.Done()
	}()

	t.handlerTaskToRunWG.Add(1)
	go func() {
		t.handlerTaskToRun(t.taskToRun, t.taskToProcess, t.stopHandlerTaskToRun)
		t.handlerTaskToRunWG.Done()
	}()

	t.handlerTaskToProcessWG.Add(1)
	go func() {
		t.handlerTaskToProcess(t.taskToProcess, t.taskDone, t.stopHandlerTaskToProcess, t.taskToSend)
		t.handlerTaskToProcessWG.Done()
	}()

	t.handlerTaskToSendWG.Add(1)
	go func() {
		t.handlerTaskToSend(t.taskToSend, t.stopHandlerTaskToSend)
		t.handlerTaskToSendWG.Done()
	}()

	t.runWorkerTaskAckWG.Add(1)
	go func() {
		t.runner.RunWorkerTaskAck(t.taskDone)
		t.runWorkerTaskAckWG.Done()
	}()

	// Todo exit all if one process is kill
	t.runWorkerTaskProviderWG.Wait()
	t.handlerTaskToRunWG.Wait()
	t.handlerTaskToProcessWG.Wait()
	t.handlerTaskToSendWG.Wait()
	t.runWorkerTaskAckWG.Wait()
	return nil
}

// StopWorker stop all goroutine and runner worker
func (t *Taskor) StopWorker() {
	t.workerStopMutex.Lock()
	defer t.workerStopMutex.Unlock()

	if !t.workerRunning {
		return
	}
	t.workerRunning = false

	// First stop consume task and wait worker stop
	log.Info("Stopping runner task provider")
	t.stopWorkerTaskProvider <- true
	t.runWorkerTaskProviderWG.Wait()

	log.Info("Stopping internal task handlers")
	t.stopHandlerTaskToRun <- true
	t.handlerTaskToRunWG.Wait()

	log.Info("Waiting last task processing")
	t.stopHandlerTaskToProcess <- true
	t.handlerTaskToProcessWG.Wait()

	log.Info("Waiting last task sending")
	t.stopHandlerTaskToSend <- true
	t.handlerTaskToSendWG.Wait()
	// Closing this chan should stop runner TaskAck
	close(t.taskDone)
	// wait last task was ACK
	log.Info("Waiting ACK last task")
	t.runWorkerTaskAckWG.Wait()
	// Stop runner
	t.runner.Stop()

	// Close other chan
	close(t.taskToRun)
	close(t.taskToProcess)
	close(t.taskToSend)
	close(t.stopWorkerTaskProvider)
	close(t.stopHandlerTaskToProcess)
	close(t.stopHandlerTaskToRun)
	close(t.stopHandlerTaskToSend)
	log.Info("Worker stopped")
}

// handlerTaskToRun handle task in chan taskToRun and process it
func (t *Taskor) handlerTaskToRun(taskToRun <-chan task.Task, taskToProcess chan<- task.Task, stop <-chan bool) {
	stopped := false
loop:
	for {
		select {
		case <-stop:
			break loop
		case queuedTask, ok := <-taskToRun:
			if !ok {
				// Chan was closed
				break loop
			}
			if queuedTask.ETA.After(time.Now()) {
				time.AfterFunc(time.Until(queuedTask.ETA), func() {
					if !stopped {
						taskToProcess <- queuedTask
					}
				})
			} else {
				// Exec task
			push:
				for {
					select {
					case <-stop:
						break loop
					case taskToProcess <- queuedTask:
						break push
					}
				}
			}
		}
	}
	stopped = true
}

func (t *Taskor) handlerTaskToSend(taskToSend <-chan task.Task, stop <-chan bool) {
loop:
	for {
		select {
		case <-stop:
			break loop
		case queuedTask, ok := <-taskToSend:
			if !ok {
				// Chan was closed
				break loop
			}
			// Retry to send the task until it works
			for {
				err := t.Send(&queuedTask)
				if err == nil {
					break
				}
				log.ErrorWithFields(fmt.Sprintf("send task error: %v", err), queuedTask)
				// We don't want to overload the runner
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// handleTaskToProcess is in charge to consume chan taskToProcess and exec task
func (t *Taskor) handlerTaskToProcess(taskToProcess <-chan task.Task, taskDone chan<- task.Task, stop <-chan bool, taskToSend chan<- task.Task) {
	// create a poll of workers to process task in concurrency
	concurrency := t.runner.GetConcurrency()
	pool := make(chan struct{}, concurrency)

	// initialize worker pool with maxWorkers workers
	go func() {
		for i := 0; i < concurrency; i++ {
			pool <- struct{}{}
		}
	}()

loop:
	for {
		select {
		case <-stop:
			break loop
		case currentTask, ok := <-taskToProcess:
			if !ok {
				// Chan was closed
				break loop
			}

			// Wait for a worker to be ready in the worker pool
			if concurrency > 0 {
				<-pool
			}

			// run task inside a go routine for parallel execution, add worker back to the pool (channel) at the end
			go func() {
				// Waiting task from runner
				err := t.execTask(&currentTask)
				// handle error (need retry/ link error / .. )
				if err != nil {
					if err == task.ErrNotRegisterd {
						if concurrency > 0 {
							pool <- struct{}{}
						}
						return
					}
					t.taskErrorHandler(&currentTask, err, taskToSend)
				} else {
					// Run child task if no error
					for _, childTask := range currentTask.ChildTasks {
						if childTask == nil {
							continue
						}
						childT := *childTask
						childT.ParentTask = &currentTask
						taskToSend <- childT
					}
					log.InfoWithFields("Task is done without error", currentTask)
				}
				// Inform runner task is finish and can be ack
				taskDone <- currentTask
				t.metric.TaskDoneWithSuccess++
				// add a worker to pool to start processing futur tasks
				if concurrency > 0 {
					pool <- struct{}{}
				}
			}()
		}
	}
}

// execTask run task function
func (t *Taskor) execTask(currentTask *task.Task) (err error) {
	Definition := t.taskList[currentTask.TaskName]
	if Definition == nil {
		log.ErrorWithFields("Task was pooled but was not register", currentTask)
		return task.ErrNotRegisterd
	}

	defer func() {
		// Handle panic in task execution, in case of panic the task is considered as in error
		if r := recover(); r != nil {
			currentTask.Error = fmt.Sprint(r)
			err = errors.New(fmt.Sprint(r))
			currentTask.DateDone = time.Now()
		}
	}()

	// Before Running task
	currentTask.DateExecuted = time.Now()
	currentTask.SetCurrentTry(currentTask.CurrentTry + 1)
	// Execute Task
	err = Definition.Run(currentTask)
	if err != nil {
		// Add error msg
		currentTask.Error = err.Error()
	}

	// After task execution
	currentTask.DateDone = time.Now()
	return err
}

// taskErrorHandler handle task error with retrying or call linked error task
func (t *Taskor) taskErrorHandler(taskToHandleError *task.Task, err error, taskToSend chan<- task.Task) {
	if err == nil {
		// task has no error to handle
		return
	}

	retry := false
	switch {
	case err == task.ErrTaskRetry:
		retry = true
	case taskToHandleError.RetryOnError:
		retry = true
	}
	// Retry if possible else call linked error task
	if retry && t.retryTaskIfPossible(taskToHandleError, taskToSend) {
		// the task has been retried
		log.InfoWithFields(fmt.Sprintf("Retry: Task failed with error: %v", err), *taskToHandleError)
		return
	}

	log.ErrorWithFields(fmt.Sprintf("Task failed with error: %v", err), *taskToHandleError)
	t.metric.TaskDoneWithError++

	// Call linked error task
	if taskToHandleError.LinkError != nil {
		// Do not use pointer here, to avoid infinite loop
		linkErrorTask := *taskToHandleError.LinkError
		linkErrorTask.ParentTask = taskToHandleError
		taskToSend <- linkErrorTask
	}
}

// retryTaskIfPossible retry task if possible return true if task is retry else false
func (t *Taskor) retryTaskIfPossible(taskToRetry *task.Task, taskToSend chan<- task.Task) bool {
	// Negative value mean infinite retry
	if taskToRetry.MaxRetry >= 0 && taskToRetry.CurrentTry > taskToRetry.MaxRetry {
		log.ErrorWithFields("Task has reached MaxRetry", taskToRetry)
		return false
	}

	// Duplicate task to avoid problem because we will repush task as a new one
	newTask := *taskToRetry
	// Adjust date when we need to retry
	newTask.ETA = time.Now().Add(newTask.CountDownRetry)
	taskToSend <- newTask
	return true
}
