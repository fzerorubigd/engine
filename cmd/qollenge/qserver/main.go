package main

import (
	"elbix.dev/engine/cmd/qollenge"
	"elbix.dev/engine/pkg/cli"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/initializer"
)

func main() {
	ctx := cli.Context()
	qollenge.InitializeConfig(ctx)

	defer initializer.Initialize(ctx)()

	grpcgw.Serve(ctx)

}
