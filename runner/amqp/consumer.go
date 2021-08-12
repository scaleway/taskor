package amqp

import (
	"time"

	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task"
	"github.com/streadway/amqp"
)

func (t *RunnerAmqp) createConsumer() <-chan amqp.Delivery {
	var msgs <-chan amqp.Delivery
	var err error
	for {
		if t.channel == nil {
			time.Sleep(errorRetryWaitTime)
			continue
		}

		msgs, err = t.channel.Consume(
			t.queueName, // queue
			"",          // consumer
			false,       // auto-ack
			false,       // exclusive
			false,       // no-local
			false,       // no-wait
			nil,         // args
		)
		if err != nil {
			time.Sleep(errorRetryWaitTime)
			continue
		}
		break
	}
	return msgs
}

// RunWorkerTaskProvider runner that consume rabbitmq and push task to taskToRun chan
func (t *RunnerAmqp) RunWorkerTaskProvider(taskToRun chan task.Task, stop <-chan bool) error {
	msgs := t.createConsumer()
loop:
	for {
		select {
		case <-stop:
			break loop
		case d, ok := <-msgs:
			if !ok {
				msgs = t.createConsumer()
			}
			// Unserialize task
			newTask := task.Task{}
			err := serializer.GetSerializer(t.serializer).Unserialize(&newTask, d.Body)
			if err != nil {
				log.Warn("[error] Cannot unserialise task, continue ...")
				continue
			}
			// Add task to mapping var
			t.addProcessingTask(newTask.RunningID, &d)

			// Send message or stop
		push:
			for {
				select {
				case taskToRun <- newTask:
					break push
				case <-stop:
					break loop
				}
			}
		}
	}
	log.Info("Consumer AMQP stopped")
	return nil
}

// RunWorkerTaskAck runner that ack message when a task is done
func (t *RunnerAmqp) RunWorkerTaskAck(taskDone <-chan task.Task) {
	for {
		taskToAck, ok := <-taskDone
		if !ok {
			// Exit if channel is close
			break
		}
		// ACK task
		delivery, err := t.getAndDeleteProcessingTask(taskToAck.RunningID)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if err := delivery.Ack(false); err != nil {
			log.InfoWithFields("Error Acking message for task", taskToAck)
			continue
		}
	}
	log.Info("Ack runner stopped")
}
