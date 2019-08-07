package user

import (
	// implementation
	_ "github.com/fzerorubigd/engine/modules/user/impl"
	// middleware
	_ "github.com/fzerorubigd/engine/modules/user/middlewares"
	// Migrations
	_ "github.com/fzerorubigd/engine/modules/user/migrations"
	// Base models and protobuf/grpc code
	_ "github.com/fzerorubigd/engine/modules/user/proto"
)
