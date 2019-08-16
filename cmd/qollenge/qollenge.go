package qollenge

import (
	"context"

	"elbix.dev/engine/modules/misc"
	"elbix.dev/engine/modules/user"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/token/jwt"
)

var (
	privateKey = config.RegisterString("secret.private", "", "RSA private key, base64 encoded")
	publicKey  = config.RegisterString("secret.public", "", "RSA public key, base64 encoded")
)

// InitializeConfig to initialize config and import packages
func InitializeConfig(ctx context.Context, initModules bool) error {
	config.Initialize(ctx, "engine", "E")

	if !initModules {
		return nil
	}
	p, err := jwt.NewJWTTokenProvider(privateKey.String(), publicKey.String())
	if err != nil {
		return err
	}

	if err := misc.SetPublicKey(publicKey.String()); err != nil {
		return err
	}
	// Its time to initialize any module
	user.SetProvider(p)
	return nil
}
