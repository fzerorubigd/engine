package main

import (
	"github.com/fzerorubigd/balloon/pkg/cli"
	"github.com/fzerorubigd/balloon/pkg/config"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/initializer"
)

func main() {
	ctx := cli.Context()

	config.Initialize("balloon", "BAL")
	defer initializer.Initialize(ctx)()

	grpcgw.Serve(ctx)

}
