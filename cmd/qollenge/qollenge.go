package qollenge

import (
	"context"

	_ "elbix.dev/engine/modules/user"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/grpcgw"
)

var (
	privateKey = config.RegisterString("secret.private", "", "RSA private key, base64 encoded")
	mods       = []func() (grpcgw.Controller, error){
		userMod,
		miscMod,
	}
)

func getPrivateKey() string {
	return privateKey.String()
}

// InitializeConfig to initialize config and import packages
func InitializeConfig(ctx context.Context, initModules bool) error {
	config.Initialize(ctx, "engine", "E")

	if !initModules {
		return nil
	}
	for i := range mods {
		m, err := mods[i]()
		if err != nil {
			return err
		}
		grpcgw.Register(m)
	}

	return nil
}
