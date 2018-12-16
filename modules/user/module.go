package user

import (
	// implementation
	_ "github.com/fzerorubigd/balloon/modules/user/impl"
	// middleware
	_ "github.com/fzerorubigd/balloon/modules/user/middlewares"
	// Migrations
	_ "github.com/fzerorubigd/balloon/modules/user/migrations"
	// Base models and protobuf/grpc code
	_ "github.com/fzerorubigd/balloon/modules/user/proto"
)
