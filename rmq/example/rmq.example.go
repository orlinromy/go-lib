package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/kelchy/go-lib/log"
	"github.com/kelchy/go-lib/rmq"
	"github.com/streadway/amqp"
)

var logger log.Log

func main() {
	logger, _ = log.New("standard")
	// Create a new RMQManager to help manage the RMQ connections
	rmqManager, rmqManagerErr := rmq.NewRMQManager(
		&rmq.QueueManagerConfig{
			Enabled:         true,
			ConnURIs:        []string{os.Getenv("RMQ_URI")},
			AutoReconnect:   true,
			EnablePublisher: true, // Enable Publisher to publish to queue
			EnableConsumer:  true, // Enable Consumer to consume from queue
			// Declares how long to wait before attempting to reconnect to RMQ upon failure
			Reconnect: rmq.Reconnect{
				Interval:   10 * time.Second,
				MaxAttempt: 2,
			},
		},
	)

	if rmqManagerErr != nil {
		logger.Error("ERR_RMQ_MANAGER_START", rmqManagerErr)
		return
	}
	logger.Out("RMQ_MANAGER_START", "RMQ Manager started")

	// Create a queue binding to an exchange
	// Exchange must be present already
	rmqQueueErr := rmqManager.CreateAndBindQueue(
		context.TODO(),
		rmq.WithQueueConfig(&rmq.QueueConfig{
			Name:       "test-queue",
			Durable:    true, // On rmq restart, the queue will be restored
			AutoDelete: true, // Once last consumer disconnects, the queue will be deleted
			Exclusive:  true, // Only accessible by the connection that declares it
			NoWait:     false,
			// Dead letter exchange and routing key to publish to (Funuctionality untested)
			Args: map[string]interface{}{
				"x-dead-letter-exchange":    "test-dlx",
				"x-dead-letter-routing-key": "event.go.retry",
			},
		}),
		rmq.WithQueueBindingConfig(&rmq.QueueBindConfig{
			Queue:      "test-queue",
			Exchange:   "test-exchange",
			BindingKey: "test-routing-key",
			NoWait:     false,
			Args:       nil,
		}))
	if rmqQueueErr != nil {
		logger.Error("ERR_RMQ-QUEUE-CONNECT", rmqQueueErr)
		return
	}

	// Functionality of event handler is declared below
	eventHandler := &Processor{}
	// Create a consumer to consume from the queue
	rmqConsumerErr := rmqManager.Consume(context.TODO(), eventHandler, rmq.WithConsumerConfig(&rmq.ConsumerConfig{
		Queue:           "test-queue",
		Enabled:         true,
		Name:            "go_queue_consumer",
		AutoAck:         false,
		Exclusive:       false,
		NoLocal:         false,
		NoWait:          false,
		Args:            nil,
		EnabledPrefetch: true,
		PrefetchCount:   1,
		PrefetchSize:    0,
		Global:          false,
	}), rmq.WithMsgRetryConfig(&rmq.MessageRetryConfig{
		Enabled:           true,
		HandleDeadMessage: true,
		RetryCountLimit:   2,
	}))
	if rmqConsumerErr != nil {
		logger.Error("ERR_RMQ_CONSUMER_CONNECT", rmqConsumerErr)
		return
	}
	logger.Out("RMQ_CONSUMER_CONNECT", "RMQ Consumer connected")

	// Create a publisher to the Exchange
	rmqPublisherErr := rmqManager.CreatePublisher(
		context.TODO(),
		rmq.NewExchange(
			&rmq.ExchangeConfig{
				Exchange:     "test-exchange",
				ExchangeType: "direct",
				Durable:      true,
				AutoDelete:   false,
				Exclusive:    false,
				NoWait:       false,
				Args:         nil,
			}, logger),
		rmq.WithPublisherConfig(&rmq.PublisherConfig{
			Enabled:               true,
			Name:                  "test-publisher",
			RoutingKey:            "test-routing-key",
			Mandatory:             false,
			Immediate:             false,
			AutoGenerateMessageID: true,
			PublisherConfirmed:    true,
			Timeout:               1,
			NoWait:                false,
		}))

	if rmqPublisherErr != nil {
		logger.Error("ERR_RMQ_PUBLISHER_CONNECT", rmqPublisherErr)
		return
	}

	logger.Out("RMQ_PUBLISHER_CONNECT", "RMQ Publisher connected")
	logger.Out("EVENT_PUBLISH", "Publishing event to RMQ")
	// Publish an event to the exchange
	publish(rmqManager, Event{Message: "test-event"})
	time.Sleep(10 * time.Second)
	publish(rmqManager, Event{Message: "test-event-after-10-seconds"})
	time.Sleep(3 * time.Second)
}

// Event processor to handle events
type Processor struct {
}

// Process the event
func (p *Processor) ProcessEvent(ctx context.Context, message rmq.IMessage) error {
	logger.Out("RMQ_EVENT_RECEIVED", string(message.Body()))
	// Do something with the message received here
	// i.e. Unmarshal and process the message
	return nil
}

// Process dead messages (Functionality untested)
func (p *Processor) ProcessDeadMessage(ctx context.Context, message rmq.IMessage, err error) error {
	logger.Error(fmt.Sprintf("ERR_RMQ_DEAD_EVENT: Headers:%s Body:%s", message.Headers(), string(message.Body())), err)
	// Do something with the dead message here
	return nil
}

// Event to be published, this event's structure should be changed to whatever you require
type Event struct {
	Message string `json:"event"`
}

func publish(r rmq.IRmqManager, e Event) {
	eventByte, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}
	_, err = r.Publish(context.Background(), eventByte,
		rmq.WithPublishingKey(amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Timestamp:    time.Now(),
		}))

	//NOTE : If you want to pass different routing key then use below code snippet
	//By default lib will use routing key from config
	// _, err = r.Publish(context.Background(), eventByte,
	// 	rmq.WithPublishingKey(amqp.Publishing{
	// 		DeliveryMode: amqp.Transient,
	// 		ContentType:  "application/json",
	// 		Timestamp:    time.Now(),
	// 	}),
	// 	rmq.WithRoutingKey("your_fav_routing_key"))

	if err != nil {
		fmt.Println(err)
	}
}
