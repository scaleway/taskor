package main

import (
	"errors"
	"log"
	"os"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/amqp"
	"github.com/scaleway/taskor/task"
)

// MyTask Task to Deploy an App
var MyTask = &task.Definition{
	Name: "MyTask",
	Run: func(task *task.Task) error {
		log.Printf("Current task")
		return errors.New("Myerror")
	},
}

// MyLinkErrorTask Task to Deploy an App
var MyLinkErrorTask = &task.Definition{
	Name: "MyLinkedErrorTask",
	Run: func(task *task.Task) error {
		log.Printf("LinkError Task id")
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
	taskManager.Handle(MyLinkErrorTask)

	if "worker" == os.Args[1] {
		taskManager.RunWorker()
	}

	if "task" == os.Args[1] {
		log.Printf("Send task")
		t, _ := task.CreateTask("MyTask", nil)
		linkTask, _ := task.CreateTask("MyLinkedErrorTask", nil)
		t.SetLinkError(linkTask)
		taskManager.Send(t)
	}

}
