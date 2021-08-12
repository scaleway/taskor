package amqp

import (
	"errors"

	"github.com/scaleway/taskor/serializer"
	"github.com/scaleway/taskor/task"
	"github.com/streadway/amqp"
)

// Send send a new task in queue
func (t *RunnerAmqp) Send(task *task.Task) error {
	var err error

	if t.channel == nil {
		return errors.New("channel is not initialized")
	}

	// Serialize Task with global serializer
	body, err := serializer.GetSerializer(t.serializer).Serialize(task)
	if err != nil {
		return err
	}

	err = t.channel.Publish(
		"",          // exchange
		t.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
