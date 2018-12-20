package kv

import (
	"time"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/redis"
)

const prefix string = "_REDIS_PREFIX_"

// TODO : maybe change it to some kind of adapter style

// StoreKey try to save the key in the key value store
func StoreKey(key, value string, d time.Duration) error {
	s := redis.Client.Set(prefix+key, value, d)
	return s.Err()
}

// MustStoreKey try to save the key value panic on error
func MustStoreKey(key, value string, d time.Duration) {
	assert.Nil(StoreKey(key, value, d))
}

// FetchKey return the key if its already there
func FetchKey(key string) (string, error) {
	a := redis.Client.Get(prefix + key)
	if err := a.Err(); err != nil {
		return "", err
	}

	return a.Val(), nil
}

// DeleteKey try to delete a key
func DeleteKey(key string) error {
	a := redis.Client.Del(prefix + key)
	return a.Err()
}

// MustDeleteKey try to delete a key panic on error
func MustDeleteKey(key string) {
	assert.Nil(DeleteKey(key))
}
