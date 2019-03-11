package main

import (
	"log"
	"os"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/amqp"
	"github.com/scaleway/taskor/task"
)

var myErrorTask = &task.Definition{
	Name: "MyTask",
	Run: func(currentTask *task.Task) error {
		log.Printf("Task try number %d", currentTask.CurrentTry)
		// Return specific error to retry
		if currentTask.CurrentTry < 20 {
			return task.ErrTaskRetry
		}
		log.Printf("Should be iteration number 20")
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

	taskManager.Handle(myErrorTask)

	if "worker" == os.Args[1] {
		taskManager.RunWorker()
	}

	if "task" == os.Args[1] {
		log.Printf("Send task")
		t, err := task.CreateTask("MyTask", nil)
		if err != nil {
			log.Fatalf(err.Error())
		}
		// Infinite retry
		t.SetMaxRetry(-1)
		taskManager.Send(t)
	}

}
