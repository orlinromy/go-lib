package rmq

import (
	"context"
)

type queueConfigKey struct{}

func WithQueueConfig(cfg *QueueConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, queueConfigKey{}, cfg)
	}
}

type queueBindingConfigKey struct{}

func WithQueueBindingConfig(cfg *QueueBindConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, queueBindingConfigKey{}, cfg)
	}
}
