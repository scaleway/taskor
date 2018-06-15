package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/amqp"
	"github.com/scaleway/taskor/task"
)

var myErrorTask = &task.Definition{
	Name: "MyTask",
	Run: func(task *task.Task) error {
		log.Printf("Task try number %d", task.CurrentTry)
		return errors.New("I am an error")
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
		t.SetMaxRetry(5).SetRetryOnError(true).SetCountDownRetry(2 * time.Second)
		taskManager.Send(t)
	}

}
