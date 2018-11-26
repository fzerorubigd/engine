package userpb

import (
	"context"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type controller struct {
}

func (uc *controller) Init(ctx context.Context, ch inprocgrpc.Channel, mux *runtime.ServeMux) {
	RegisterHandlerUserSystem(&ch, NewUserController())
	cl := NewUserSystemChannelClient(&ch)

	assert.Nil(RegisterUserSystemHandlerClient(ctx, mux, cl))
}

func init() {
	grpcgw.Register(&controller{})
}
