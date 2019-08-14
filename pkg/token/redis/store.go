package redis

import (
	"encoding/json"
	"time"

	"elbix.dev/engine/pkg/token"

	"elbix.dev/engine/pkg/kv"
	"elbix.dev/engine/pkg/random"
)

type redisStore struct {
}

func (redisStore) Store(data map[string]interface{}, exp time.Duration) (string, error) {
	str, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	key := <-random.ID
	if err := kv.StoreKey(key, string(str), exp); err != nil {
		return "", err
	}

	return key, nil
}

func (redisStore) Fetch(token string) (map[string]interface{}, error) {
	res, err := kv.FetchKey(token)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(res), &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// NewRedisTokenProvider return a new provider
func NewRedisTokenProvider() token.Provider {
	return &redisStore{}
}
