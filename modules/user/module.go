package user

import (
	"elbix.dev/engine/modules/user/impl"
	"elbix.dev/engine/modules/user/middlewares"
	"elbix.dev/engine/pkg/token"

	// Migrations
	_ "elbix.dev/engine/modules/user/migrations"
	// Base models and protobuf/grpc code
	_ "elbix.dev/engine/modules/user/proto"
)

// TODO : some sort of dependency injection in module level

// SetProvider for setting the token provider
func SetProvider(p token.Provider) {
	impl.SetProvider(p)
	middlewares.SetProvider(p)
}
