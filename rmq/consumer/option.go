package consumer

import (
	"context"
	"reflect"
)

type Options struct {
	// for alternative data
	Context        context.Context
	optionalParams []KV
}

type KV struct {
	Key   interface{}
	Value interface{}
}

// pass cloption to config struct
type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type OptionalKV func() (interface{}, interface{})

func Opt(key, value interface{}) OptionalKV {
	if key == nil || !reflect.TypeOf(key).Comparable() {
		return nil
	}
	return func() (interface{}, interface{}) {
		return key, value
	}

}
func NewOptionsFromOptional(opts ...OptionalKV) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		if o == nil {
			continue
		}
		key, val := o()

		options.optionalParams = append(options.optionalParams, KV{key, val})
		options.Context = context.WithValue(options.Context, key, val)
	}

	return options
}

func (o Options) GetValue(key interface{}) interface{} {
	for _, kv := range o.optionalParams {
		if kv.Key == key {
			return kv.Value
		}
	}
	return nil
}

func (o Options) GetValues() []KV {
	return o.optionalParams
}

type noWaitKey struct{}

func WithNoWait(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, noWaitKey{}, a)
	}
}

type consumerConfigKey struct{}

func WithConsumerConfig(cfg *ConsumerConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, consumerConfigKey{}, cfg)
	}
}

type msgRetryConfigKey struct{}

func WithMsgRetryConfig(cfg *MessageRetryConfig) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, msgRetryConfigKey{}, cfg)
	}
}

type multipleKey struct{}

func WithMultiple(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, multipleKey{}, a)
	}
}

type requeueKey struct{}

func WithRequeue(a bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, requeueKey{}, a)
	}
}
