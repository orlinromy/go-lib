package consumer

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type IEventHandler interface {
	HandleEvent(context.Context, IMessage)
	Retry(context.Context, IMessage, error)
	HandleDeadMessage(context.Context, IMessage, error)
}

type IClientHandler interface {
	ProcessEvent(context.Context, IMessage) error
	ProcessDeadMessage(context.Context, IMessage, error) error
}

type EventHandler struct {
	processor   IClientHandler
	logger      ILogger
	retryConfig *MessageRetryConfig
}

func NewEventHandler(processor IClientHandler, logger ILogger, retryConfig *MessageRetryConfig) IEventHandler {
	return &EventHandler{
		processor:   processor,
		logger:      logger,
		retryConfig: retryConfig,
	}
}

func (e *EventHandler) HandleEvent(ctx context.Context, message IMessage) {
	err := e.processor.ProcessEvent(ctx, message)
	if err != nil {
		if e.retryConfig.Enabled {
			e.Retry(ctx, message, err)
			return
		}
		e.HandleDeadMessage(ctx, message, err)
		return
	}
}

func (e *EventHandler) Retry(ctx context.Context, message IMessage, err error) {
	headers := message.Headers()
	if headers == nil { // in case of 1st retry no headers are present
		errAck := message.Ack(false, WithRequeue(false))
		if errAck != nil {
			e.logger.Error(fmt.Sprintf("ERR_EVENT_HANDLER-FAIL-MSG-ACK-%s", message.GetID()), errAck)
		}
		return
	}
	if xDeathContent, ok := headers["x-death"].([]interface{}); ok {
		for _, content := range xDeathContent {
			table, _ := content.(amqp.Table)
			retryCount, _ := table["count"].(int64)
			if int(retryCount) <= e.retryConfig.RetryCountLimit {
				errAck := message.Ack(false, WithRequeue(false))
				if errAck != nil {
					e.logger.Error(fmt.Sprintf("ERR_EVENT_HANDLER-FAIL-MSG-ACK-%s", message.GetID()), errAck)
					return
				}
				return
			}
			e.HandleDeadMessage(ctx, message, fmt.Errorf("retry count %d exceeded for msg %s", int(retryCount), message.GetID()))
		}
	}
	e.HandleDeadMessage(ctx, message, err)
}

func (e *EventHandler) HandleDeadMessage(ctx context.Context, message IMessage, err error) {
	if e.retryConfig.HandleDeadMessage {
		handleDMErr := e.processor.ProcessDeadMessage(ctx, message, err)
		if handleDMErr != nil {
			e.logger.Error(fmt.Sprintf("ERR_EVENT_HANDLER-DEAD-MSG-FAIL-MSG-ACK-%s", message.GetID()), handleDMErr)
		}
	}
}
