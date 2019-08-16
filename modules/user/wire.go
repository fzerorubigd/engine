// +build wireinject

package user

import (
	"elbix.dev/engine/modules/user/impl"
	// middlewares
	_ "elbix.dev/engine/modules/user/middlewares"
	// migrations
	_ "elbix.dev/engine/modules/user/migrations"
	userpb "elbix.dev/engine/modules/user/proto"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/sec"
	"elbix.dev/engine/pkg/token/jwt"
	"github.com/google/wire"
)

// UserSet is the builder used to build this module
var UserSet = wire.NewSet(
	wire.Bind(new(grpcgw.Controller), new(userpb.WrappedUserSystemController)),
	sec.ParseRSAPrivateKeyFromBase64PEM,
	jwt.NewJWTTokenProvider,
	userpb.NewWrappedUserSystemServer,
	impl.NewUserController,
)
