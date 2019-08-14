package user

import (
	userpb "elbix.dev/engine/modules/user/proto"
	"elbix.dev/engine/pkg/token"

	// Migrations
	_ "elbix.dev/engine/modules/user/migrations"
	// Base models and protobuf/grpc code
	_ "elbix.dev/engine/modules/user/proto"
)

// TODO : some sort of dependency injection in module level

// SetProvider for setting the token provider
func SetProvider(p token.Provider) {
	userpb.SetProvider(p)
}
