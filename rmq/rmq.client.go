package rmq

import (
	"context"
	"fmt"
	"time"

	"github.com/kelchy/go-lib/log"
	"github.com/streadway/amqp"
)

const (
	DefaultReconnectBackOffTime = time.Millisecond * 1000
)

type IRmqManager interface {
	CreateExchange(context.Context, ...Option) (IExchange, error)
	CreatePublisher(context.Context, IExchange, ...Option) error
	CreateAndBindQueue(context.Context, ...Option) error
	Publish(context.Context, []byte, ...Option) (string, error)
	Consume(context.Context, IClientHandler, ...Option) error
}

type rmqManager struct {
	opts   Options
	config *QueueManagerConfig
	logger log.Log

	conPool      IConnectionPool
	conn         *amqp.Connection
	errorConn    chan *amqp.Error
	errorPubChan chan *amqp.Error
	errorConChan chan *amqp.Error
	pubChan      *amqp.Channel
	conChan      *amqp.Channel

	exchange  IExchange
	publisher IPublisher
	consumer  IConsumer
}

func NewRMQManager(config *QueueManagerConfig, opts ...Option) (IRmqManager, error) {
	options := NewOptions(opts...)

	if !config.Enabled {
		return nil, queueManagerIsDisabledError
	}

	if config.Reconnect.Interval == 0 {
		config.Reconnect.Interval = DefaultReconnectBackOffTime
	}

	rmqManager := &rmqManager{
		opts:   options,
		config: config,
	}

	var Log, _ = log.New("standard")
	rmqManager.logger = Log

	rmqManager.conPool = newConnectionPool(rmqManager.config.ConnURIs...)

	rmqManager.connect(false)
	go rmqManager.reconnect()

	go rmqManager.listenOnChanClose()

	rmqManager.logger.Debug("RMQ_MANAGER", "rmq manager created successfully")
	return rmqManager, nil
}

func (r *rmqManager) openChannel() (*amqp.Channel, error) {
	if r.conn == nil || r.conn.IsClosed() {
		return nil, connIsNotOpened
	}
	return r.conn.Channel()
}

func (r *rmqManager) reconnect() {
	for {
		err := <-r.errorConn
		if err != nil {
			r.logger.Out("RMQ_MANAGER-RECONNECT", "Reconnecting to RMQ")
			r.connect(true)
		}
	}
}

func (r *rmqManager) listenOnChanClose() {
	for {
		select {
		case err := <-r.errorPubChan:
			if err != nil {
				r.logger.Error("ERR_RMQ_PUBLISHER-PUB-CHAN-ERR", err)
				if !r.conn.IsClosed() {
					errClose := r.conn.Close()
					if errClose != nil {
						r.logger.Error("ERR_RMQ_PUBLISHER-PUB-CHAN-ERR-IN-CONN-CLOSE", errClose)
					}
				}
			}
		case err := <-r.errorConChan:
			if err != nil {
				r.logger.Error("ERR_RMQ_CONSUMER-CON-CHAN-ERR", err)
				if !r.conn.IsClosed() {
					errClose := r.conn.Close()
					if errClose != nil {
						r.logger.Error("ERR_RMQ_CONSUMER-CON-CHAN-ERR-IN-CONN-CLOSE", errClose)
					}
				}
			}
		}
	}
}

func (r *rmqManager) connect(reconnect bool) {
	for {
		r.logger.Debug("RMQ_MANAGER", "connecting to RMQ")

		conn, err := r.conPool.GetCon()
		if err == nil {
			r.conn = conn
			r.errorConn = make(chan *amqp.Error)
			r.conn.NotifyClose(r.errorConn)

			if r.config.EnablePublisher {
				c, err := r.openChannel()
				if err != nil {
					r.logger.Error("ERR_RMQ_MANAGER-FAIL-OPEN-PUB-CHANNEL", err)
					return
				}
				r.pubChan = c
				r.errorPubChan = make(chan *amqp.Error)
				r.pubChan.NotifyClose(r.errorPubChan)
			}

			if r.config.EnableConsumer {
				c, err := r.openChannel()
				if err != nil {
					r.logger.Error("ERR_RMQ_MANAGER-FAIL-OPEN-CON-CHANNEL", err)
					return
				}
				r.conChan = c
				r.errorConChan = make(chan *amqp.Error)
				r.conChan.NotifyClose(r.errorConChan)
				if r.consumer != nil {
					_ = r.consumer.Start(context.Background(), r.conChan)
				}
			}

			r.logger.Out("RMQ_MANAGER", "connected to RMQ")
			return
		}

		r.logger.Error("ERR_RMQ_MANAGER-FAIL-CONNECT", err)
		time.Sleep(r.config.Reconnect.Interval)
	}
}

func (r *rmqManager) CreateExchange(ctx context.Context, opts ...Option) (IExchange, error) {
	options := NewOptions(opts...)

	exchangeConfig, ok := options.Context.Value(exchangeConfigKey{}).(*ExchangeConfig)
	if !ok || exchangeConfig == nil {
		return nil, invalidExchangeConfig
	}

	exchange := NewExchange(exchangeConfig, r.logger)
	err := exchange.Create(r.pubChan)
	if err != nil {
		return nil, err
	}
	r.exchange = exchange
	return exchange, nil
}

func (r *rmqManager) CreatePublisher(ctx context.Context, ex IExchange, opts ...Option) error {
	options := NewOptions(opts...)

	publisherConfig, ok := options.Context.Value(publisherConfigKey{}).(*PublisherConfig)
	if !ok || publisherConfig == nil {
		return invalidPublisherConfig
	}

	publisher := NewPublisher(r.logger, publisherConfig, ex, r.pubChan)
	r.publisher = publisher
	return nil
}

func (r *rmqManager) Publish(ctx context.Context, data []byte, opts ...Option) (string, error) {
	options := NewOptions(opts...)

	pubConfig := r.publisher.GetConfig()
	publishing, ok := options.Context.Value(publishingKey{}).(amqp.Publishing)
	if !ok {
		return "", invalidPublishArgs
	}
	autoGenerateMessageID, ok := options.Context.Value(autoGenerateMessageID{}).(bool)
	if !ok {
		autoGenerateMessageID = pubConfig.AutoGenerateMessageID
	}

	mandatory, ok := options.Context.Value(mandatoryKey{}).(bool)
	if !ok {
		mandatory = pubConfig.Mandatory
	}

	immediate, ok := options.Context.Value(immediateKey{}).(bool)
	if !ok {
		mandatory = pubConfig.Immediate
	}

	if autoGenerateMessageID {
		uuid, err := NewUUID()
		if err != nil {
			return "", IDGenerationFailed
		}
		publishing.MessageId = uuid
	}

	exchange, ok := options.Context.Value(exchangeKey{}).(IExchange)
	if !ok {
		exchange = r.publisher.GetExchange()
	}

	exchangeConfig := exchange.GetConfig()
	publishing.Body = data

	routingKey, ok := options.Context.Value(routingKey{}).(string)
	if !ok {
		//This will use routingKey from config is client has not provided in func args
		routingKey = pubConfig.RoutingKey
	}

	err := r.pubChan.Publish(exchangeConfig.Exchange, routingKey, mandatory, immediate, publishing)
	if err != nil {
		r.logger.Error("ERR_RMQ_MANAGER-FAIL-PUBLISH", err)
		return "", err
	}
	r.logger.Debug("RMQ_MANAGER", fmt.Sprintf("msg published to exchange => %s with id => %s", exchangeConfig.Exchange, publishing.MessageId))
	return publishing.MessageId, nil
}

func (r *rmqManager) CreateAndBindQueue(ctx context.Context, opts ...Option) error {
	options := NewOptions(opts...)

	qConfig, ok := options.Context.Value(queueConfigKey{}).(*QueueConfig)
	if !ok || qConfig == nil {
		return invalidQueueConfig
	}

	qBindConfig, ok := options.Context.Value(queueBindingConfigKey{}).(*QueueBindConfig)
	if !ok || qBindConfig == nil {
		return invalidQueueBindConfig
	}

	que := NewQueue(qConfig, r.logger)
	err := que.Create(r.conChan)
	if err != nil {
		return err
	}

	return que.Bind(r.conChan, qBindConfig)
}

func (r *rmqManager) Consume(ctx context.Context, clientHandler IClientHandler, opts ...Option) error {
	options := NewOptions(opts...)

	config, ok := options.Context.Value(consumerConfigKey{}).(*ConsumerConfig)
	if !ok || config == nil {
		return invalidConsumerConfig
	}

	retryConfig, ok := options.Context.Value(msgRetryConfigKey{}).(*MessageRetryConfig)
	if !ok || retryConfig == nil {
		retryConfig = new(MessageRetryConfig)
		retryConfig.Enabled = false
	}

	handler := NewEventHandler(clientHandler, r.logger, retryConfig)
	r.consumer = NewConsumer(config, r.logger, handler)
	return r.consumer.Start(ctx, r.conChan)
}
