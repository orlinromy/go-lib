package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	rabbitmq "github.com/kelchy/go-lib/rmq/consumer"
)

func main() {
	conn, err := rabbitmq.NewConn(
		"amqps://hpyoucun:Yec9xQqm8ZFizmZshyqjuELwqrkT79ng@armadillo.rmq.cloudamqp.com/hpyoucun",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Verify that exchange exists (Not needed if it is declared in NewConsumer)
	if err := rabbitmq.DeclareExchange(
		conn,
		rabbitmq.ExchangeOptions{
			Name:       "events",
			Kind:       "direct",
			Declare:    true,
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		},
	); err != nil {
		fmt.Println(err)
		return
	}

	// Verify that queue exists (Not needed if it is declared in NewConsumer)
	if err := rabbitmq.DeclareQueue(
		conn,
		rabbitmq.QueueOptions{
			Name:       "my_queue",
			Declare:    true,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
	); err != nil {
		fmt.Println(err)
		return
	}

	// Verify that queue is bound to exchange (Not needed binding is declared in NewConsumer)
	if err := rabbitmq.DeclareBinding(
		conn,
		rabbitmq.BindingDeclareOptions{
			QueueName:    "my_queue",
			ExchangeName: "events",
			RoutingKey:   "test_routing_key_1",
			NoWait:       false,
			Args:         nil,
			Declare:      true,
		}); err != nil {
		fmt.Println(err)
		return
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed 1: %s, %v", string(d.MessageId), string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"my_queue",
		rabbitmq.WithConsumerOptionsConcurrency(2),
		rabbitmq.WithConsumerOptionsRoutingKey("test_routing_key_1"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
	)
	if err != nil {
		fmt.Println(err)
	}
	defer consumer.Close()

	consumer2, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed 2: %s, %v", string(d.MessageId), string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"my_queue",
		rabbitmq.WithConsumerOptionsConcurrency(2),
		rabbitmq.WithConsumerOptionsRoutingKey("test_routing_key_2"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer2.Close()

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("stopping consumer")
}
