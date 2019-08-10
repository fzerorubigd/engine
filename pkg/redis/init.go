package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"elbix.dev/engine/pkg/assert"
	"elbix.dev/engine/pkg/initializer"
	"elbix.dev/engine/pkg/log"
)

var (
	// Client the actual pool to use with redis
	client *redis.Pool
	all    []initializer.Simple
	lock   sync.RWMutex
)

type initRedis struct {
}

// Initialize try to create a redis pool
func (i *initRedis) Initialize(ctx context.Context) {
	client = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				network.String(),
				fmt.Sprintf("%s:%d", address.String(), port.Int()),
				redis.DialDatabase(db.Int()),
				redis.DialPassword(password.String()),
			)
		},
		IdleTimeout:     time.Minute,
		MaxActive:       poolsize.Int(),
		MaxIdle:         5,
		Wait:            true,
		MaxConnLifetime: 10 * time.Minute,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > 10*time.Minute {
				return errors.New("old connection, refreshing")
			}
			_, err := c.Do("PING")
			return err
		},
	}

	for i := range all {
		all[i].Initialize()
	}
	log.Debug("redis is ready.")

	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(client.Close())
		log.Debug("redis finalized.")
	}()
}

// Connection return a new connection from the pool
func Connection() redis.Conn {
	return client.Get()
}

// Register a new object to inform it after redis is loaded
func Register(in initializer.Simple) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, in)
}

func init() {
	// Redis must be before mysql so the cache work on queries
	initializer.Register(&initRedis{}, -1)
}
