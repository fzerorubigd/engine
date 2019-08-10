package user

import (
	// implementation
	_ "elbix.dev/engine/modules/user/impl"
	// middleware
	_ "elbix.dev/engine/modules/user/middlewares"
	// Migrations
	_ "elbix.dev/engine/modules/user/migrations"
	// Base models and protobuf/grpc code
	_ "elbix.dev/engine/modules/user/proto"
)
