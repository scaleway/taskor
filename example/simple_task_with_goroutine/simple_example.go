package main

import (
	"log"
	"time"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/goroutine"
	"github.com/scaleway/taskor/task"
)

// MyTaskParameter parameter for task MyTask
type MyTaskParameter struct {
	MyParameter string
}

// MyTask Task to Deploy an App
var MyTask = &task.Definition{
	Name: "MyTask",
	Run: func(task *task.Task) error {

		// Get Task Param
		var param MyTaskParameter
		if err := task.UnserializeParameter(&param); err != nil {
			return err
		}
		log.Printf("Task id %s", task.RunningID)
		log.Printf("Task queued at %s", task.DateQueued)
		log.Printf("Task executed at %s", task.DateExecuted)
		log.Printf("With paramter %s", param.MyParameter)
		return nil
	},
}

func main() {

	config := goroutine.RunnerConfig{
		MaxBufferedMessage: 0,
	}
	taskManager, err := taskor.New(goroutine.New(config))
	if err != nil {
		log.Fatalf(err.Error())
	}
	taskManager.Handle(MyTask)
	// You need to run worker with Runner
	go taskManager.RunWorker()

	log.Printf("Send task")
	t, _ := task.CreateTask("MyTask", MyTaskParameter{MyParameter: "simple"})
	taskManager.Send(t)

	time.Sleep(1 * time.Second)
	taskManager.StopWorker()

	// Display metrics
	metric := taskManager.GetMetrics()
	log.Printf("Task done with error %d", metric.TaskDoneWithError)
	log.Printf("Task done without error %d", metric.TaskDoneWithSuccess)
	log.Printf("Task sent %d", metric.TaskSent)

}
