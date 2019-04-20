package redis

import (
	"github.com/fzerorubigd/balloon/pkg/config"
)

var (
	address  = config.RegisterString("services.redis.address", "127.0.0.1", "redis address host")
	network  = config.RegisterString("services.redis.network", "tcp", "redis network")
	port     = config.RegisterInt("services.redis.port", 6379, "redis port")
	password = config.RegisterString("services.redis.password", "", "redis password")
	poolsize = config.RegisterInt("services.redis.poolsize", 20, "redis pool size")
	db       = config.RegisterInt("services.redis.db", 1, "redis db number")
)
