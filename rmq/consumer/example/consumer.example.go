package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kelchy/go-lib/rmq/consumer"
)

func main() {
	eventProcessor := &EventProcessor{}
	err := consumer.New(
		consumer.DefaultConnectionConfig([]string{os.Getenv("RMQ_URI")}),
		// Queue names should be the same in QueueConfig and QueueBindConfig
		consumer.DefaultQueueConfig("test-queue-logging"),
		consumer.DefaultQueueBindConfig("test-exchange", "test-queue-logging", "test-routing-key"),
		consumer.DefaultConsumerConfig("test-consumer"),
		consumer.DefaultMessageRetryConfig(),
		eventProcessor,
		consumer.DefaultLogger())
	if err != nil {
		fmt.Println("failed to create consumer: ", err)
	}

	// Leave the consumer running for 30 seconds before exiting, only for example purposes
	time.Sleep(30 * time.Second)
}

type EventProcessor struct{}

func (*EventProcessor) ProcessEvent(ctx context.Context, message consumer.IMessage) error {
	fmt.Printf("Recieved message: ID: %s, Message: %s\n", message.GetID(), string(message.Body()))
	return nil
}

func (*EventProcessor) ProcessDeadMessage(ctx context.Context, message consumer.IMessage, err error) error {
	fmt.Printf("Recieved dead message: ID: %s, Message: %s, Error: %v", message.GetID(), string(message.Body()), err)
	return nil
}
