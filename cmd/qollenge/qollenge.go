package qollenge

import (
	"context"

	"elbix.dev/engine/pkg/config"
	// misc module
	_ "elbix.dev/engine/modules/misc"
	// user module
	_ "elbix.dev/engine/modules/user"
)

// InitializeConfig to initializ config and import packages
func InitializeConfig(ctx context.Context) {
	config.Initialize(ctx, "engine", "E")
}
