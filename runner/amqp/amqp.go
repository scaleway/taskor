package amqp

import (
	"fmt"
	"sync"
	"time"

	"github.com/scaleway/taskor/log"
	"github.com/scaleway/taskor/serializer"
	"github.com/streadway/amqp"
)

// Time to wait before retry after queue error
var errorRetryWaitTime = 1 * time.Second

// RunnerAmqpConfig config use for amqp runner
type RunnerAmqpConfig struct {
	// AmqpUrl url of rabbitmq (ex: amqp://guest:guest@localhost:5672/)
	AmqpURL      string
	ExchangeName string
	QueueName    string
	QueueDurable bool
	Concurrency  int
}

// NewConfig return a new RunnerAmqpConfig with default value
func NewConfig() RunnerAmqpConfig {
	config := RunnerAmqpConfig{
		AmqpURL:      "amqp://guest:guest@localhost:5672/",
		QueueName:    "taskor_queue",
		QueueDurable: false,
		Concurrency:  1,
	}
	return config
}

// RunnerAmqp struct
type RunnerAmqp struct {
	amqpURL     string
	queueName    string
	queueDurable bool
	concurrency  int
	serializer  serializer.Type

	// Amqp element
	conn             *amqp.Connection
	channel          *amqp.Channel
	rabbitCloseError chan *amqp.Error

	// Map between taskId and message
	processingTask      map[string]*amqp.Delivery
	mutexProcessingTask sync.Mutex
}

// New create a new runner
func New(amqpConfig RunnerAmqpConfig) *RunnerAmqp {
	runner := &RunnerAmqp{}
	runner.amqpURL = amqpConfig.AmqpURL
	runner.queueName = amqpConfig.QueueName
	runner.queueDurable = amqpConfig.QueueDurable
	runner.serializer = serializer.TypeJSON
	runner.concurrency = amqpConfig.Concurrency
	return runner
}

// GetConcurrency - get concurrency configuration
func (t *RunnerAmqp) GetConcurrency() int {
	return t.concurrency
}

// Init connection
func (t *RunnerAmqp) Init() error {
	// Init proccessing mapping between task and message
	// This is used to ack message
	t.processingTask = make(map[string]*amqp.Delivery)

	// Connect to RabbitMQ
	err := t.amqpConnect()
	if err != nil {
		// If amqp not ready do not block connection retry
		go t.amqpRetryConnect()
	}
	return nil
}

// Stop Close channel & connection
func (t *RunnerAmqp) Stop() error {
	t.channel.Close()
	t.conn.Close()
	return nil
}

func (t *RunnerAmqp) amqpConnect() error {
	log.Info("Connection to RabbitMQ")
	var err error

	conn, err := amqp.Dial(t.amqpURL)
	if err != nil {
		return err
	}
	t.conn = conn

	channel, err := t.conn.Channel()
	if err != nil {
		return err
	}
	t.channel = channel

	err = t.prepareQueue()
	if err != nil {
		return err
	}

	// Go routine to handle connection failure
	t.rabbitCloseError = make(chan *amqp.Error)
	t.conn.NotifyClose(t.rabbitCloseError)
	go t.handleAMQPFailure()

	log.Info("RabbitMq connection OK")
	return nil
}

// amqpRetryConnect infinite loop trying to connect to amqp, break when connected
func (t *RunnerAmqp) amqpRetryConnect() {
	for {
		err := t.amqpConnect()
		if err != nil {
			log.Error("Error on rabbitmq connection: " + err.Error())
			time.Sleep(errorRetryWaitTime)
			continue
		}
		break
	}
}

// handleAMQPFailure handle AMQP disconnection
func (t *RunnerAmqp) handleAMQPFailure() {

	// Wait for a Close notification
	rabbitErr := <-t.rabbitCloseError
	if rabbitErr != nil {
		log.Error("Received disconnection event")
		t.amqpRetryConnect()
	}

}

func (t *RunnerAmqp) prepareQueue() error {
	_, err := t.channel.QueueDeclare(
		t.queueName,    // name
		t.queueDurable, // queueDurable
		false,          // delete when usused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *RunnerAmqp) addProcessingTask(taskRunningID string, d *amqp.Delivery) {
	t.mutexProcessingTask.Lock()
	defer t.mutexProcessingTask.Unlock()

	t.processingTask[taskRunningID] = d
}

func (t *RunnerAmqp) getAndDeleteProcessingTask(taskRunningID string) (*amqp.Delivery, error) {
	t.mutexProcessingTask.Lock()
	defer t.mutexProcessingTask.Unlock()

	d := t.processingTask[taskRunningID]
	if d == nil {
		return nil, fmt.Errorf("[error]Processing task unreachable : %s", taskRunningID)
	}

	delete(t.processingTask, taskRunningID)
	return d, nil
}
