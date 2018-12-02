package main

import (
	"github.com/fzerorubigd/balloon/modules"
	"github.com/fzerorubigd/balloon/pkg/cli"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/initializer"
)

func main() {
	ctx := cli.Context()
	modules.InitializeConfig()

	defer initializer.Initialize(ctx)()

	grpcgw.Serve(ctx)

}
