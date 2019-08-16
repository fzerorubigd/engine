// +build wireinject

package misc

import (
	"elbix.dev/engine/modules/misc/impl"
	miscpb "elbix.dev/engine/modules/misc/proto"
	_ "elbix.dev/engine/modules/user/middlewares"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/sec"
	"github.com/google/wire"
)

// MiscSet is the builder used to build this module
var MiscSet = wire.NewSet(
	wire.Bind(new(grpcgw.Controller), new(miscpb.WrappedMiscSystemController)),
	sec.ParseRSAPrivateKeyFromBase64PEM,
	sec.ExtractPublicFromPrivate,
	miscpb.NewWrappedMiscSystemServer,
	impl.NewMiscController,
)
