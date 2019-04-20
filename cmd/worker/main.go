package main

import (
	"flag"

	"github.com/fzerorubigd/balloon/cmd"
	"github.com/fzerorubigd/balloon/pkg/cli"
	"github.com/fzerorubigd/balloon/pkg/initializer"
	"github.com/fzerorubigd/balloon/pkg/job"
	"github.com/fzerorubigd/chapar/workers"
)

func main() {
	flag.Parse()

	ctx := cli.Context()
	cmd.InitializeConfig()

	defer initializer.Initialize(ctx)()

	job.Process(ctx, workers.WithParallelLimit(10), workers.WithRetryCount(1))
}
