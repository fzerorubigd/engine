package cmd

import (
	"context"

	"github.com/fzerorubigd/engine/pkg/config"
	// misc module
	_ "github.com/fzerorubigd/engine/modules/misc"
	// user module
	_ "github.com/fzerorubigd/engine/modules/user"
)

// InitializeConfig to initializ config and import packages
func InitializeConfig(ctx context.Context) {
	config.Initialize(ctx, "engine", "E")
}
