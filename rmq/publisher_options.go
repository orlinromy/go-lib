package rmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

type publisherConfigKey struct{}

func WithPublisherConfig(cfg *PublisherConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, publisherConfigKey{}, cfg)
	}
}

type autoGenerateMessageID struct{}

func WithAutoGenerateMessageID(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, autoGenerateMessageID{}, a)
	}
}

type mandatoryKey struct{}

func WithMandatory(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, mandatoryKey{}, a)
	}
}

type immediateKey struct{}

func WithImmediate(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, immediateKey{}, a)
	}
}

type routingKey struct{}

func WithRoutingKey(a string) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, routingKey{}, a)
	}
}

type publishingKey struct{}

func WithPublishingKey(a amqp.Publishing) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, publishingKey{}, a)
	}
}

type timeoutKey struct{}

// in second
func WithTimeout(a time.Duration) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, timeoutKey{}, a)
	}
}

type exchangeKey struct{}

// in second
func WithExchange(a IExchange) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, exchangeKey{}, a)
	}
}
