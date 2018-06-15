package main

import (
	"log"
	"os"

	"github.com/scaleway/taskor"
	"github.com/scaleway/taskor/runner/amqp"
	"github.com/scaleway/taskor/task"
)

var parentTask = &task.Definition{
	Name: "ParentTask",
	Run: func(task *task.Task) error {
		log.Printf("Parent Task")
		return nil
	},
}

var childTask1 = &task.Definition{
	Name: "ChildTask1",
	Run: func(task *task.Task) error {
		log.Printf("Child Task 1")
		return nil
	},
}

var childTask2 = &task.Definition{
	Name: "ChildTask2",
	Run: func(task *task.Task) error {
		log.Printf("Child Task 2")
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
	taskManager.Handle(parentTask)
	taskManager.Handle(childTask1)
	taskManager.Handle(childTask2)

	if "worker" == os.Args[1] {
		taskManager.RunWorker()
	}

	if "task" == os.Args[1] {
		// This will run ParentTask -> (ChildTask1 | ChildTask2)
		ex1pTask, _ := task.CreateTask("ParentTask", nil)
		ex1cTask1, _ := task.CreateTask("ChildTask1", nil)
		ex1cTask2, _ := task.CreateTask("ChildTask2", nil)
		ex1pTask.AddChild(ex1cTask1).AddChild(ex1cTask2)
		taskManager.Send(ex1pTask)

		// This will run ParentTask -> ChildTask1 -> ChildTask2
		ex2pTask, _ := task.CreateTask("ParentTask", nil)
		ex2cTask1, _ := task.CreateTask("ChildTask1", nil)
		ex2cTask2, _ := task.CreateTask("ChildTask2", nil)
		ex2pTask.AddChild(ex2cTask1.AddChild(ex2cTask2))
		taskManager.Send(ex2pTask)
	}

}
