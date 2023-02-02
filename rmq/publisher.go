package rmq

import (
	Logger "github.com/kelchy/go-lib/log"

	"github.com/streadway/amqp"
)

type IPublisher interface {
	GetConfig() *PublisherConfig
	GetExchange() IExchange
}

type Publisher struct {
	logger   Logger.Log
	config   *PublisherConfig
	exchange IExchange
}

func NewPublisher(l Logger.Log, c *PublisherConfig, ex IExchange, ch *amqp.Channel) IPublisher {
	return &Publisher{
		logger:   l,
		config:   c,
		exchange: ex,
	}
}

func (p *Publisher) GetConfig() *PublisherConfig {
	return p.config
}

func (p *Publisher) GetExchange() IExchange {
	return p.exchange
}
