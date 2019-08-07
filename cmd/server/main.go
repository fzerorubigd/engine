package main

import (
	"github.com/fzerorubigd/engine/cmd"
	"github.com/fzerorubigd/engine/pkg/cli"
	"github.com/fzerorubigd/engine/pkg/grpcgw"
	"github.com/fzerorubigd/engine/pkg/initializer"
)

func main() {
	ctx := cli.Context()
	cmd.InitializeConfig(ctx)

	defer initializer.Initialize(ctx)()

	grpcgw.Serve(ctx)

}
