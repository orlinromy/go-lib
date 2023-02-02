package rmq

import (
	"context"
)

type exchangeConfigKey struct{}

func WithExchangeConfig(cfg *ExchangeConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, exchangeConfigKey{}, cfg)
	}
}
