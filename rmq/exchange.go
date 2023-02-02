package rmq

import (
	Logger "github.com/kelchy/go-lib/log"
	"github.com/streadway/amqp"
)

type IExchange interface {
	Name() string
	GetConfig() *ExchangeConfig
	Create(*amqp.Channel) error
}

type Exchange struct {
	config *ExchangeConfig
	logger Logger.Log
}

func NewExchange(config *ExchangeConfig, logger Logger.Log) IExchange {
	return &Exchange{
		config: config,
		logger: logger,
	}
}

func (e *Exchange) Name() string {
	return e.config.Exchange
}

func (e *Exchange) GetConfig() *ExchangeConfig {
	return e.config
}

func (e *Exchange) Create(channel *amqp.Channel) error {
	c := e.config
	return channel.ExchangeDeclare(c.Exchange, c.ExchangeType, c.Durable, c.AutoDelete, c.Exclusive, c.NoWait, c.Args)
}
