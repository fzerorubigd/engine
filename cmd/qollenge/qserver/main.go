package main

import (
	"github.com/fzerorubigd/engine/cmd/qollenge"
	"github.com/fzerorubigd/engine/pkg/cli"
	"github.com/fzerorubigd/engine/pkg/grpcgw"
	"github.com/fzerorubigd/engine/pkg/initializer"
)

func main() {
	ctx := cli.Context()
	qollenge.InitializeConfig(ctx)

	defer initializer.Initialize(ctx)()

	grpcgw.Serve(ctx)

}
