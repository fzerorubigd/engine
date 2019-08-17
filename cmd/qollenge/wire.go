// +build wireinject

package qollenge

import (
	"elbix.dev/engine/modules/misc"
	"elbix.dev/engine/modules/user"
	"elbix.dev/engine/pkg/grpcgw"
	"github.com/google/wire"
)

func userMod() (grpcgw.Controller, error) {
	panic(wire.Build(getPrivateKey, user.UserSet))
}

func miscMod() (grpcgw.Controller, error) {
	panic(wire.Build(getPrivateKey, misc.MiscSet))
}

