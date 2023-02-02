package rmq

import (
	"fmt"

	Logger "github.com/kelchy/go-lib/log"
	"github.com/streadway/amqp"
)

type IQueue interface {
	GetName() string
	Create(*amqp.Channel) error
	Bind(*amqp.Channel, *QueueBindConfig) error
}

type Queue struct {
	config *QueueConfig
	logger Logger.Log
}

func NewQueue(config *QueueConfig, logger Logger.Log) IQueue {
	return &Queue{
		config: config,
		logger: logger,
	}
}

func (q *Queue) GetName() string {
	return q.config.Name
}

func (q *Queue) Create(channel *amqp.Channel) error {
	que, err := channel.QueueDeclare(q.config.Name, q.config.Durable, q.config.AutoDelete, q.config.Exclusive, q.config.NoWait, q.config.Args)
	if err != nil {
		return err
	}
	q.logger.Debug("RMQ_QUEUE", fmt.Sprintf("create queue => %s success", que.Name))
	return nil
}

func (q *Queue) Bind(channel *amqp.Channel, config *QueueBindConfig) error {
	err := channel.QueueBind(config.Queue, config.BindingKey, config.Exchange, config.NoWait, config.Args)
	if err != nil {
		return err
	}
	q.logger.Debug("RMQ_QUEUE", fmt.Sprintf("bind queue => %s with exchange => %s success", config.Queue, config.Exchange))
	return nil
}
