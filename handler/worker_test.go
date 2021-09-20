package handler

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	runnerMock "github.com/scaleway/taskor/runner/mock"
	"github.com/scaleway/taskor/task"
)

func TestTaskor_retryTaskIfPossible(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().Send(gomock.Any()).AnyTimes()
	ta, _ := New(mockRunner)

	taskToSend := make(chan task.Task, 100)

	tests := []struct {
		name        string
		taskToRetry *task.Task
		want        bool
	}{
		{
			name: "task need retry",
			want: true,
			taskToRetry: &task.Task{
				CurrentTry: 1,
				MaxRetry:   2,
			},
		},
		{
			name: "task infitine retry",
			want: true,
			taskToRetry: &task.Task{
				CurrentTry: 10,
				MaxRetry:   -1,
			},
		},
		{
			name: "task maxretry",
			want: false,
			taskToRetry: &task.Task{
				CurrentTry: 3,
				MaxRetry:   2,
			},
		},
		{
			name: "task equal maxretry",
			want: true,
			taskToRetry: &task.Task{
				CurrentTry: 2,
				MaxRetry:   2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ta.retryTaskIfPossible(tt.taskToRetry, taskToSend); got != tt.want {
				t.Errorf("Taskor.retryTaskIfPossible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskor_execTask(t *testing.T) {
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().Send(gomock.Any()).AnyTimes()

	ta, _ := New(mockRunner)
	// Create ok task
	errorTest := errors.New("error task return")
	var taskTest = task.Definition{
		Name: "test",
		Run:  func(t *task.Task) error { return errorTest },
	}
	testTask, _ := task.CreateTask("test", nil)

	// register Definition for "test"
	ta.Handle(&taskTest)

	t.Run("execTaskNotRegistered", func(t *testing.T) {
		testTaskNotRegister, _ := task.CreateTask("testNotRegister", nil)
		err := ta.execTask(testTaskNotRegister)
		if err != task.ErrNotRegisterd {
			t.Errorf("Task does not return errorTest")
		}
	})

	t.Run("execTask", func(t *testing.T) {
		err := ta.execTask(testTask)

		if err != errorTest {
			t.Errorf("Task does not return errorTest")
		}

		if testTask.ID == "" {
			t.Errorf("Task ID is nil")
		}

		if testTask.CurrentTry != 1 {
			t.Errorf("Task SetCurrentTry is not 1 : %d", testTask.CurrentTry)
		}

		if time.Time.IsZero(testTask.DateExecuted) {
			t.Errorf("Task DateExecuted is nil")
		}

		if time.Time.IsZero(testTask.DateDone) {
			t.Errorf("Task DateDone is nil")
		}
	})
}

func TestTaskor_execTaskPanic(t *testing.T) {
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().Send(gomock.Any()).AnyTimes()

	ta, _ := New(mockRunner)
	// Create ok task
	errorTest := errors.New("unexpected error")
	var taskTest = task.Definition{
		Name: "test",
		Run:  func(t *task.Task) error { panic("unexpected error") },
	}
	testTask, _ := task.CreateTask("test", nil)

	// register Definition for "test"
	ta.Handle(&taskTest)

	t.Run("execTask", func(t *testing.T) {
		err := ta.execTask(testTask)

		if err.Error() != errorTest.Error() {
			t.Errorf("Task does not return errorTest")
		}

		if testTask.ID == "" {
			t.Errorf("Task ID is nil")
		}

		if testTask.CurrentTry != 1 {
			t.Errorf("Task SetCurrentTry is not 1 : %d", testTask.CurrentTry)
		}

		if time.Time.IsZero(testTask.DateExecuted) {
			t.Errorf("Task DateExecuted is nil")
		}

		if time.Time.IsZero(testTask.DateDone) {
			t.Errorf("Task DateDone is nil")
		}
	})
}

func TestTaskor_taskErrorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()

	t.Run("task no error", func(t *testing.T) {
		ta, _ := New(mockRunner)
		taskToSend := make(chan task.Task, 100)
		testTask, _ := task.CreateTask("test", nil)
		testTask.ID = "testtaskid"
		ta.taskErrorHandler(testTask, nil, taskToSend)
		if ta.metric.TaskDoneWithError != 0 {
			t.Errorf("Metric is incremented")
		}
	})

	t.Run("task retry error", func(t *testing.T) {
		ta, _ := New(mockRunner)
		taskToSend := make(chan task.Task, 100)
		testTask, _ := task.CreateTask("test", nil)
		testTask.ID = "testtaskid"

		testTask.MaxRetry = -1
		testTask.RetryOnError = false
		ta.taskErrorHandler(testTask, task.ErrTaskRetry, taskToSend)

		sentTask := <-taskToSend
		// Retried task should keep taskID
		if sentTask.ID != testTask.ID {
			t.Errorf("Wrong task ID: %s", sentTask.ID)
		}
		if ta.metric.TaskDoneWithError != 0 {
			t.Errorf("Metric is incremented")
		}
	})

	t.Run("task retryOnerror", func(t *testing.T) {
		ta, _ := New(mockRunner)
		taskToSend := make(chan task.Task, 100)
		testTask, _ := task.CreateTask("test", nil)
		testTask.ID = "testtaskid"

		testTask.MaxRetry = -1
		testTask.RetryOnError = true
		ta.taskErrorHandler(testTask, errors.New("task custom error"), taskToSend)

		sentTask := <-taskToSend
		// Retried task should keep taskID
		if sentTask.ID != testTask.ID {
			t.Errorf("Wrong task ID: %s", sentTask.ID)
		}
		if ta.metric.TaskDoneWithError != 0 {
			t.Errorf("Metric is incremented")
		}
	})

	t.Run("task retryOnerror without error task", func(t *testing.T) {
		ta, _ := New(mockRunner)
		taskToSend := make(chan task.Task, 100)
		testTask, _ := task.CreateTask("test", nil)
		testTask.ID = "testtaskid"
		testTask.MaxRetry = -1
		testTask.RetryOnError = false
		ta.taskErrorHandler(testTask, errors.New("task custom error"), taskToSend)
		if ta.metric.TaskDoneWithError != 1 {
			t.Errorf("Metric is not incremented")
		}
	})

	t.Run("task no retry with error task", func(t *testing.T) {
		ta, _ := New(mockRunner)
		taskToSend := make(chan task.Task, 100)
		testTask, _ := task.CreateTask("test", nil)
		testTask.ID = "testtaskid"
		errorTask, _ := task.CreateTask("linkedErrorTask", nil)
		testTask.SetLinkError(errorTask)
		testTask.RetryOnError = false
		testTask.MaxRetry = 0
		ta.taskErrorHandler(testTask, errors.New("task custom error"), taskToSend)
		sentTask := <-taskToSend
		if sentTask.TaskName != "linkedErrorTask" {
			t.Errorf("Wrong task name: %s", sentTask.TaskName)
		}
		if sentTask.ParentTask.TaskName != "test" {
			t.Errorf("Wrong parent task name: %s", sentTask.TaskName)
		}
		if ta.metric.TaskDoneWithError != 1 {
			t.Errorf("Metric is not incremented")
		}
	})
}

func TestTaskor_handlerTaskToProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().GetConcurrency().Return(1).AnyTimes()
	ta, _ := New(mockRunner)

	var taskTest = task.Definition{
		Name: "test",
		Run:  func(t *task.Task) error { return nil },
	}

	// register Definition for "test"
	ta.Handle(&taskTest)

	testTask, _ := task.CreateTask("test", nil)
	testTask.ID = "testtaskid"
	taskToProcess := make(chan task.Task)
	taskToSend := make(chan task.Task)
	taskDone := make(chan task.Task, 100)
	stop := make(chan bool, 1)

	t.Run("stop handlerTaskToProcess", func(t *testing.T) {
		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		stop <- true
		ta.handlerTaskToProcess(taskToProcess, taskDone, stop, taskToSend)
		timer.Stop()
	})

	t.Run("handlerTaskToProcess", func(t *testing.T) {
		go func() {
			ta.handlerTaskToProcess(taskToProcess, taskDone, stop, taskToSend)
		}()
		// Insert a task to exec
		taskToProcess <- *testTask
		processTask := <-taskDone
		if processTask.RunningID != testTask.RunningID {
			t.Errorf("Wrong task to process")
		}
		// stop the goroutine function
		stop <- true
	})

	t.Run("task with children", func(t *testing.T) {
		child1Task, _ := task.CreateTask("test", nil)
		child2Task, _ := task.CreateTask("test", nil)
		testTask.AddChild(child1Task)
		testTask.AddChild(child2Task)

		go ta.handlerTaskToProcess(taskToProcess, taskDone, stop, taskToSend)
		taskToProcess <- *testTask
		child1ToSend := <-taskToSend
		child2ToSend := <-taskToSend
		if child1ToSend.ParentTask.TaskName != "test" {
			t.Errorf("Wrong parent task name: %s", child1ToSend.TaskName)
		}
		if child2ToSend.ParentTask.TaskName != "test" {
			t.Errorf("Wrong parent task name: %s", child1ToSend.TaskName)
		}
		if child1ToSend.ID != child1Task.ID || child2ToSend.ID != child2Task.ID {
			t.Errorf("Wrong task was sent")
		}
		<-taskDone
		stop <- true
	})

	t.Run("close chan handlerTaskToProcess", func(t *testing.T) {
		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		close(taskToProcess)
		ta.handlerTaskToProcess(taskToProcess, taskDone, stop, taskToSend)
		timer.Stop()
	})
}

func TestTaskor_handlerTaskToSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	ta, _ := New(mockRunner)

	var taskTest = task.Definition{
		Name: "test",
		Run:  func(t *task.Task) error { return nil },
	}

	// register Definition for "test"
	ta.Handle(&taskTest)

	testTask, _ := task.CreateTask("test", nil)
	testTask.ID = "testtaskid"
	taskToSend := make(chan task.Task)
	stop := make(chan bool, 1)

	t.Run("stop handlerTaskToSend", func(t *testing.T) {
		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		stop <- true
		ta.handlerTaskToSend(taskToSend, stop)
		timer.Stop()
	})

	t.Run("handlerTaskToSend", func(t *testing.T) {
		lambda := func(sendTask *task.Task) {
			if sendTask.ID != testTask.ID {
				panic("Wrong task was sent")
			}
		}
		mockRunner.EXPECT().Send(gomock.Any()).Do(lambda)
		go func() {
			ta.handlerTaskToSend(taskToSend, stop)
		}()
		// Insert a task to send
		taskToSend <- *testTask
		// stop the goroutine function
		stop <- true
	})

	t.Run("sending fails twice", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		mockRunner.EXPECT().Send(gomock.Any()).Return(errors.New("some error")).Times(2)
		mockRunner.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
		go func() {
			defer wg.Done()
			ta.handlerTaskToSend(taskToSend, stop)
		}()
		// Insert a task to send
		taskToSend <- *testTask
		// stop the goroutine function
		stop <- true

		wg.Wait()
	})

	t.Run("close chan handlerTaskToSend", func(t *testing.T) {
		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		close(taskToSend)
		ta.handlerTaskToSend(taskToSend, stop)
		timer.Stop()
	})
}

func TestTaskor_handlerTaskToRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	ta, _ := New(mockRunner)

	testTask, _ := task.CreateTask("test", nil)
	testTask.ID = "testtaskid"

	t.Run("stop handlerTaskToRun", func(t *testing.T) {
		taskToProcess := make(chan task.Task)
		taskToRun := make(chan task.Task)
		stop := make(chan bool, 1)

		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		stop <- true
		ta.handlerTaskToRun(taskToRun, taskToProcess, stop)
		timer.Stop()
	})

	t.Run("handlerTaskToRun", func(t *testing.T) {
		taskToProcess := make(chan task.Task)
		taskToRun := make(chan task.Task)
		stop := make(chan bool, 1)

		go func() {
			ta.handlerTaskToRun(taskToRun, taskToProcess, stop)
		}()
		// Insert a task to run
		taskToRun <- *testTask
		processTask := <-taskToProcess
		if processTask.RunningID != testTask.RunningID {
			t.Errorf("Wrong task to process")
		}
		// stop the goroutine function
		stop <- true
	})

	t.Run("handlerTaskToRun with ETA in future", func(t *testing.T) {
		taskToProcess := make(chan task.Task)
		taskToRun := make(chan task.Task)
		stop := make(chan bool, 1)

		go func() {
			ta.handlerTaskToRun(taskToRun, taskToProcess, stop)
		}()
		// Insert a task to run
		testTask.ETA = time.Now().Add(10 * time.Second)
		taskToRun <- *testTask

		// I don't know how to do that in other way
		// PR or help are accepted :D
		time.Sleep(100 * time.Millisecond)
		select {
		case <-taskToProcess:
			t.Errorf("Task was processed")
		default:
		}

		// stop the goroutine function
		stop <- true
	})

	t.Run("close chan handlerTaskToRun", func(t *testing.T) {
		taskToProcess := make(chan task.Task)
		taskToRun := make(chan task.Task)
		stop := make(chan bool, 1)

		timer := time.AfterFunc(1*time.Second, func() {
			panic("Process don't stop")
		})
		close(taskToRun)
		ta.handlerTaskToRun(taskToRun, taskToProcess, stop)
		timer.Stop()
	})
}

func TestTaskor_StartStopWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()
	mockRunner.EXPECT().GetConcurrency().Return(1).AnyTimes()
	mockRunner.EXPECT().RunWorkerTaskProvider(gomock.Any(), gomock.Any()).
		Do(func(taskToRun chan task.Task, stop <-chan bool) {
			<-stop
		})
	mockRunner.EXPECT().RunWorkerTaskAck(gomock.Any()).
		Do(func(taskDone <-chan task.Task) {
			<-taskDone
		})
	mockRunner.EXPECT().Stop()

	ta, _ := New(mockRunner)
	go ta.RunWorker()
	timer := time.AfterFunc(8*time.Second, func() {
		panic("Process don't stop")
	})

	// Todo remove this
	time.Sleep(1 * time.Second)
	ta.StopWorker()
	timer.Stop()
}

func TestTaskor_StartWorkerAlreadyStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRunner := runnerMock.NewMockRunner(ctrl)
	mockRunner.EXPECT().Init().AnyTimes()

	ta, _ := New(mockRunner)
	ta.workerRunning = true
	err := ta.RunWorker()
	if err != errorWorkerAlreadyRunning {
		t.Errorf("worker can be start twice, err %v", err)
	}
}
