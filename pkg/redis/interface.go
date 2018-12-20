package redis

import (
	"github.com/go-redis/redis"
)

// RedisClient is the interface to use instead of normal redis client and cover both redis and redis-cluster
type RedisClient interface {
	redis.Cmdable
	Close() error
}
