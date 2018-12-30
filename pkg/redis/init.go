package redis

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/initializer"
	"github.com/fzerorubigd/balloon/pkg/log"
	"github.com/go-redis/redis"
)

var (
	// Client the actual pool to use with redis
	Client ClientInterface
	all    []initializer.Simple
	lock   sync.RWMutex
)

type initRedis struct {
}

// Initialize try to create a redis pool
func (i *initRedis) Initialize(ctx context.Context) {
	if cluster.Bool() {
		endpoints, err := lookup(address.String())
		assert.Nil(err)
		for i := range endpoints {
			endpoints[i] = fmt.Sprintf("%s:%d", endpoints[i], port.Int())
		}
		Client = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:    endpoints,
				Password: password.String(),
				PoolSize: poolsize.Int(),
			},
		)
	} else {
		Client = redis.NewClient(
			&redis.Options{
				Network:  "tcp",
				Addr:     fmt.Sprintf("%s:%d", address.String(), port.Int()),
				Password: password.String(),
				PoolSize: poolsize.Int(),
				DB:       db.Int(),
			},
		)
	}
	// PING the server to make sure every thing is fine
	if err := Client.Ping().Err(); err != nil {
		log.Fatal("Can not connect to redis", log.Err(err))
	}

	for i := range all {
		all[i].Initialize()
	}
	log.Debug("redis is ready.")
	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(Client.Close())
		log.Debug("redis finalized.")
	}()
}

// Register a new object to inform it after redis is loaded
func Register(in initializer.Simple) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, in)
}

func lookup(svcName string) ([]string, error) {
	var endpoints []string
	_, srvRecords, err := net.LookupSRV("", "", svcName)
	if err != nil {
		return endpoints, err
	}
	for _, srvRecord := range srvRecords {
		// The SRV records ends in a "." for the root domain
		ep := fmt.Sprintf("%v", srvRecord.Target[:len(srvRecord.Target)-1])
		endpoints = append(endpoints, ep)
	}
	fmt.Print(endpoints)
	return endpoints, nil
}

func init() {
	// Redis must be before mysql so the cache work on queries
	initializer.Register(&initRedis{}, -1)
}
