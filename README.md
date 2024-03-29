# TASKOR
Task queue over RabbitMQ lib system.


## Installing
It is easy to use Taskor.
In first time, use `go get` to clone the latest version of the library.

```
go get -u github.com/scaleway/taskor
```

Next, include Taskor in your go file:

```go
import "github.com/scaleway/taskor"
```


## How it works

### Define a task
``` go
type MyTaskParameter struct {
	MyParameter        string
}

var MyTask = &taskor.Definition{
	Name: "MyTask",
	Run: func(task *task.Task) error {

		// Get Task Param
		var param MyTaskParameter
		if err := task.UnserializeParameter(&param); err != nil {
			return err
		}
		log.Printf("With paramter %s", param.MyParameter)
		return nil
	},
}
```

### Send a task
``` go
taskManager := taskor.New(runner.RunnerTypeAmqp, config)
taskManager.Handle(MyTask)
MyTask, _ := t.CreateTask("MyTask", param)
t.Send(MyTask)
```

### Running worker
``` go
// With AMQP driver
config := amqp.NewConfig()
// Feel free to update your configuration
config.AmqpURL = "amqp://guest:guest@localhost:5672/"
config.ExchangeName = "myexchange"
config.QueueName = "taskor_queue"
config.Concurrency = 5
amqpRunner := amqp.New(config)
taskManager := taskor.New(amqpRunner)
taskManager.Handle(MyTask)
taskManager.RunWorker()
```

### Example
See files in example directory

## Advanced features

### Concurrency

By default, each worker can only handle one task at the same time (wait for the current task to finish processing before processing the next one). With `Concurrency` property on your runner configuration, you can set a maximum number of workers processing tasks concurrently.

Taskor uses goroutines to run multiple tasks in parallel. To define max number of workers processing tasks at the same time (done at worker initialization):
``` go
amqpConfig := amqp.NewConfig()
// set Concurrency on runner configuration
amqpConfig.Concurrency = 5

amqpRunner := amqp.New(amqpConfig)
taskManager := taskor.New(amqpRunner)
// ...
taskManager.RunWorker()
```

### Retry
To define MaxRetry allowed for a task:
``` go
MyTask.MaxRetry = 5
// or you can do this
MyTask.SetMaxRetry(5)
```

Task can be retry if an error is return:
``` go
MyTask.RetryOnError = true
// or you can do this
MyTask.SetRetryOnError(true)
```

You can chain both
``` go
MyTask.SetRetryOnError(true).SetMaxRetry(5)
```

If you don't want to retry on each error but in a specific case, you just had to return task.ErrTaskRetry as error of your task.
``` go
var MyErrorTask = &task.Definition{
	Name: "MyTask",
	Run: func(currentTask *task.Task) error {
		log.Printf("Task try number %d", currentTask.CurrentTry)
		// Return specific error to retry
		if currentTask.CurrentTry < 20 {
			return task.ErrTaskRetry # This will perform a retry
		}
		log.Printf("Should be iteration number 20")
		return nil
	},
}
```

To sum up:

* Taskor doesn't retry a task by default (`MaxRetry: 0`).
* Taskor retries only if:
  * `MaxRetry` is defined (use `-1` for infinite retries),
  * a task returns `task.ErrTaskRetry`,
  * a task returns any error when `RetryOnError` is `true`.

### LinkError
LinkError is used to link a task that will be run when a task ending whith error and can't be retry.

To link a task to an other:
``` go
MyTask.SetLinkError(myLinkErrorTask)
```

### Chain
Chain task is possible adding a child task to an other task. Child will be execute only if parent task is successful
``` go
MyTask.AddChild(MyOtherTask)
```

### ParentTask
In case of LinkError or ChildTask, it's possible to access to parent information using task.ParentTask
``` go
 func(task *task.Task) error {
		log.Printf("MyParent error was %s", task.ParentTask.Error)
		return nil
 }
```

### Define a custom logger
A taskor logger should implement this interface:
``` go
type Logger interface {
	Debug(msg string, extraFields map[string]interface{})
	Info(msg string, extraFields map[string]interface{})
	Warn(msg string, extraFields map[string]interface{})
	Error(msg string, extraFields map[string]interface{})
}
```
To change taskor logger :
``` go
import 	taskorLogger "github.com/scaleway/taskor/log"

taskorLogger.SetLogger(&LogrusTaskor{})
```
See in example dir to see how to implement logrus

### Get some metrics

``` go
metric := taskManager.GetMetrics()
log.Printf("Task done with error %d", metric.TaskDoneWithError)
log.Printf("Task done without error %d", metric.TaskDoneWithSuccess)
log.Printf("Task sent %d", metric.TaskSent)
```

# Other links:
* [HowItWorks](doc/HowItWorks.md)

# Dev

## Necessaries packages

```
go get github.com/golang/mock/gomock
go get github.com/streadway/amqp
```
## Run Test
```
make test
```

## Build
```
make build
```

## Generate mock
```
make mock
```

## Import mock

Taskor is mocked in mock, generated by mockgen (https://github.com/golang/mock)

If you want to mock taskor in your project, you can use this example:

``` go
import "github.com/scaleway/taskor/mock"

mockTaskor := mock_taskor.NewMockTaskManager(ctrl)
mockTaskor.EXPECT().Send(gomock.Any()).Times(1).Do(func(taskToRun *taskorTask.Task){
    ...
})
```
