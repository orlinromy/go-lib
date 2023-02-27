package redis

import (
	"context"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// Get - implementation of redis GET, ctx can be nil
func (r Client) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		r.log.Error("REDIS_GET", err)
		return val, err
	}
	return val, nil
}

// Keys - implementation of redis KEYS, ctx can be nil
func (r Client) Keys(ctx context.Context, match string) ([]string, error) {
	var list []string
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.Client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			r.log.Error("REDIS_KEYS", err)
			return list, err
		}
		list = append(list, keys...)
		if cursor == 0 {
			break
		}
	}
	return list, nil
}

// TTL - implementation of redis TTL, ctx can be nil
// - returns in time.Duration which is always in nanoseconds, so we accept a params to convert in library
func (r Client) TTL(ctx context.Context, key string, precision string) (time.Duration, error) {
	val, err := r.Client.TTL(ctx, key).Result()
	if err != nil && err != redis.Nil {
		r.log.Error("REDIS_TTL", err)
		return 0, err
	}

	// handle conversion
	switch strings.ToLower(precision) {
	case "hours":
		val = val / time.Hour
	case "minutes":
		val = val / time.Minute
	case "seconds":
		val = val / time.Second
	case "microseconds":
		val = val / time.Microsecond
	case "millisecond":
		val = val / time.Millisecond
	default:
		val = val / time.Second
	}

	return val, err
}
