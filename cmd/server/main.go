package main

import (
	"net/http"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/modules/user/impl"
	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/fraugster/cli"
)

func main() {
	ctx := cli.Context()

	var c inprocgrpc.Channel
	userpb.RegisterHandlerUserSystem(&c, userimpl.NewUserController())

	cl := userpb.NewUserSystemChannelClient(&c)

	mux := runtime.NewServeMux()
	if err := userpb.RegisterUserSystemHandlerClient(ctx, mux, cl); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
