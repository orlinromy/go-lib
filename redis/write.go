package redis

import (
	"time"
	"errors"
	redis "github.com/go-redis/redis/v8"
)

var lockPrefix = "lock_"

// Lock - implementation of redis distributed lock
func (r Client) Lock(key string, ttl time.Duration) (bool, error) {
	tmpkey := lockPrefix + key
	args := redis.SetArgs{
		TTL:	ttl,
		Mode:	"NX",
	}
	s := r.Client.SetArgs(r.ctx, tmpkey, "", args)
	res, err := s.Result()

	// if key is locked, returns redis.Nil
	// implied that res != OK
	if err == redis.Nil {
		return false, nil
	}
	if res == "OK" {
		return true, nil
	}
	if err != nil {
		r.log.Error("REDIS_LOCK", err)
		return false, err
	}
	// by default don't give lock
	return false, nil
}

// Unlock - implementation of redis distributed unlock
func (r Client) Unlock(key string) (bool, error) {
	tmpkey := lockPrefix + key
	res, err := r.Del(tmpkey)
	if err != nil {
		err = errors.New("Unlock: failed")
		r.log.Error("REDIS_UNLOCK", err)
		return false, err
	}
	resbool := false
	if res > 0 {
		resbool = true
	}
	return resbool, err
}

// Set - implementation of redis SET
func (r Client) Set(key string, value string, ttl time.Duration) (string, error) {
	res, err := r.Client.Set(r.ctx, key, value, ttl).Result()
	if err != nil {
		r.log.Error("REDIS_SET", err)
		return res, err
	}
	return res, nil
}

// SetNX - implementation of redis SET with NX flag
func (r Client) SetNX(key string, value string, ttl time.Duration) (bool, error) {
	res, err := r.Client.SetNX(r.ctx, key, value, ttl).Result()
	if err != nil {
		r.log.Error("REDIS_SETNX", err)
		return res, err
	}
	return res, nil
}

// Del - implementation of redis DEL
func (r Client) Del(key string) (int64, error) {
	res, err := r.Client.Del(r.ctx, key).Result()
	if err != nil {
		r.log.Error("REDIS_DEL", err)
		return res, err
	}
	return res, nil
}
