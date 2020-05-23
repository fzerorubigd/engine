package main

import (
	"syscall"

	"elbix.dev/engine/app/common"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/initializer"
	"elbix.dev/engine/pkg/log"
	"github.com/fzerorubigd/clictx"
)

func main() {
	ctx := clictx.Context(
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	if err := common.InitializeConfig(ctx, true); err != nil {
		log.Fatal("Dependency injection failed", log.Err(err))
	}

	defer initializer.Initialize(ctx)()

	if err := grpcgw.Serve(ctx); err != nil {
		log.Error("Serve failed with an error", log.Err(err))
	}
}
