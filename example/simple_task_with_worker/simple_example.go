package main

import (
	"log"
	"os"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/amqp"
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
	if len(os.Args) != 2 {
		log.Printf("Please chose worker or task")
		return
	}

	config := amqp.NewConfig()
	taskManager, err := taskor.New(amqp.New(config))
	if err != nil {
		log.Fatalf(err.Error())
	}
	taskManager.Handle(MyTask)

	if "worker" == os.Args[1] {
		taskManager.RunWorker()
	}

	if "task" == os.Args[1] {
		log.Printf("Send task")

		// Quick way to run a task
		t, _ := task.CreateTask("MyTask", MyTaskParameter{MyParameter: "simple"})
		taskManager.Send(t)
	}

}
