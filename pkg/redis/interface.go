package redis

import (
	"github.com/go-redis/redis"
)

// ClientInterface is the interface to use instead of normal redis client and cover both redis and redis-cluster
type ClientInterface interface {
	redis.Cmdable
	Close() error
}
