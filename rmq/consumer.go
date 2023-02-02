package rmq

import (
	"context"
	"fmt"

	"github.com/kelchy/go-lib/log"
	"github.com/streadway/amqp"
)

type IConsumer interface {
	Start(context.Context, *amqp.Channel) error
}

type Consumer struct {
	config  *ConsumerConfig
	logger  log.Log
	handler IEventHandler
}

func NewConsumer(config *ConsumerConfig, logger log.Log, handler IEventHandler) IConsumer {
	return &Consumer{
		config:  config,
		logger:  logger,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context, channel *amqp.Channel) error {
	if c.config.EnabledPrefetch {
		err := channel.Qos(c.config.PrefetchCount, c.config.PrefetchSize, c.config.Global)
		if err != nil {
			c.logger.Error("ERR_CONSUMER-FAILED-PREFETCH-ADD", err)
			return err
		}
	}
	msgs, err := channel.Consume(c.config.Queue, c.config.Name, c.config.AutoAck,
		c.config.Exclusive, c.config.NoLocal, c.config.NoWait, c.config.Args)
	if err != nil {
		c.logger.Error(fmt.Sprintf("ERR_CONSUMER-FAILED-CONSUMER-REGISTER-%s", c.config.Name), err)
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		for msg := range msgs {
			c.handler.HandleEvent(ctx, NewMessage(msg))
		}
	}(msgs)

	return nil
}
