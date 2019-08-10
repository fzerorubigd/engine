package kv

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"elbix.dev/engine/pkg/assert"
	aredis "elbix.dev/engine/pkg/redis"
)

const prefix string = "_REDIS_PREFIX_"

// TODO : maybe change it to some kind of adapter style

// StoreKey try to save the key in the key value store
func StoreKey(key, value string, d time.Duration) error {
	conn := aredis.Connection()
	defer func() {
		_ = conn.Close()
	}()

	_, err := redis.String(conn.Do("PSETEX", prefix+key, d.Nanoseconds()/1000, value))
	return errors.Wrap(err, "set failed")
}

// MustStoreKey try to save the key value panic on error
func MustStoreKey(key, value string, d time.Duration) {
	assert.Nil(StoreKey(key, value, d))
}

// FetchKey return the key if its already there
func FetchKey(key string) (string, error) {
	conn := aredis.Connection()
	defer func() {
		_ = conn.Close()
	}()

	s, err := redis.String(conn.Do("GET", prefix+key))
	if err != nil {
		return "", errors.Wrap(err, "get failed")
	}
	return s, nil
}

// DeleteKey try to delete a key
func DeleteKey(key string) error {
	conn := aredis.Connection()
	defer func() {
		_ = conn.Close()
	}()

	_, err := conn.Do("DEL", prefix+key)
	return errors.Wrap(err, "can not delete the key")
}

// MustDeleteKey try to delete a key panic on error
func MustDeleteKey(key string) {
	assert.Nil(DeleteKey(key))
}

// TTLKey return the ttl of a key
func TTLKey(key string) (time.Duration, error) {
	conn := aredis.Connection()
	defer func() {
		_ = conn.Close()
	}()

	ttl, err := redis.Int64(conn.Do("TTL", prefix+key))
	assert.Nil(err)
	if ttl < 0 {
		return 0, errors.New("key not found or had no ttl")
	}

	return time.Duration(1000 * ttl), nil
}

// MustTTLKey is the must version of the ttl key func
func MustTTLKey(key string) time.Duration {
	ttl, err := TTLKey(key)
	assert.Nil(err)
	return ttl
}
