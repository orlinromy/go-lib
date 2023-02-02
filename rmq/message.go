package rmq

import (
	"github.com/streadway/amqp"
)

type IMessage interface {
	GetID() string
	Ack(flag bool, opts ...Option) error
	Headers() map[string]interface{}
	Body() []byte
}

type Message struct {
	delivery amqp.Delivery
}

func NewMessage(d amqp.Delivery) IMessage {
	return &Message{
		delivery: d,
	}
}

func (m *Message) GetID() string {
	return m.delivery.MessageId
}

func (m *Message) Ack(flag bool, opts ...Option) error {
	options := NewOptions(opts...)
	multiple, _ := options.Context.Value(multipleKey{}).(bool)
	requeue, _ := options.Context.Value(requeueKey{}).(bool)
	if flag {
		err := m.delivery.Ack(multiple)
		if err != nil {
			return err
		}
	} else {
		err := m.delivery.Reject(requeue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) Body() []byte {
	return m.delivery.Body
}

func (m *Message) Headers() map[string]interface{} {
	return m.delivery.Headers
}
