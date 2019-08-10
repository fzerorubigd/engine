package redis

import (
	"context"

	"github.com/fzerorubigd/chapar/drivers/redis"
	"github.com/fzerorubigd/chapar/middlewares/storage"
	"github.com/fzerorubigd/chapar/tasks"
	"github.com/fzerorubigd/chapar/workers"

	"elbix.dev/engine/pkg/assert"
)

type redisStorage struct {
	prefix string
}

// Store is an example of store. mostly you do not want to return err that easily
// also maybe better storage (like database?) or ...
func (s *redisStorage) Store(ctx context.Context, task *tasks.Task, e error) (es error) {
	d, err := task.Marshal()
	if err != nil {
		return err // maybe just log?
	}
	key := s.prefix + task.ID.String()
	c := Connection()
	defer func() {
		_ = c.Close()
	}()

	_, err = c.Do("HSET", key, "TASK", string(d))
	if err != nil {
		return err
	}
	if e != nil {
		_, err = c.Do("HSET", key, "ERR", e.Error())
		if err != nil {
			return err
		}
		_, err = c.Do("HINCRBY", key, "REDELIVER", "1")
		if err != nil {
			return err
		}
	}

	// lets set the time for 3 days to not bloat the server
	_, err = c.Do("EXPIRE", key, 72*60*60)
	return err
}

// NewJobStore return a job store in redis
func NewJobStore(prefix string) storage.Store {
	return &redisStorage{prefix: prefix}
}

// NewDriver returns new worker driver using redis
func NewDriver(ctx context.Context, prefix string) workers.Driver {
	driver, err := redis.NewDriver(
		ctx,
		redis.WithQueuePrefix(prefix),
		redis.WithRedisPool(client),
	)
	assert.Nil(err)
	return driver
}
