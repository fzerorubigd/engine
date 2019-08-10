package postgres

import (
	"elbix.dev/engine/pkg/config"
)

var (
	// retryMax  = config.RegisterDuration("services.postgres.max_retry_connection", time.Minute, "max time app should fallback to get mysql connection")
	user    = config.RegisterString("services.postgres.user", "engine", "postgres user")
	dbname  = config.RegisterString("services.postgres.db", "engine", "postgres database")
	pass    = config.RegisterString("services.postgres.password", "bita123", "postgres password")
	host    = config.RegisterString("services.postgres.host", "localhost", "postgres host")
	port    = config.RegisterInt("services.postgres.port", 5432, "postgres port")
	maxIdle = config.RegisterInt("services.postgres.max_idle_connection", 10, "postgres maximum idle connection")
	maxCon  = config.RegisterInt("services.postgres.max_connection", 150, "postgres  maximum connection")
	// develMode = config.RegisterBoolean("core.devel_mode", false, "development mode?")
	sslmode = config.RegisterString("services.postgres.sslmode", "disable", "sslmode for postgres")
)
