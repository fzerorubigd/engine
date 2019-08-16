package qollenge

import (
	"context"

	"elbix.dev/engine/modules/misc"
	"elbix.dev/engine/modules/user"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/sec"
	"elbix.dev/engine/pkg/token/jwt"
	"github.com/pkg/errors"
)

var (
	privateKey = config.RegisterString("secret.private", "", "RSA private key, base64 encoded")
)

// InitializeConfig to initialize config and import packages
func InitializeConfig(ctx context.Context, initModules bool) error {
	config.Initialize(ctx, "engine", "E")

	if !initModules {
		return nil
	}
	key, err := sec.ParseRSAPrivateKeyFromBase64PEM(privateKey.String())
	if err != nil {
		return errors.Wrap(err, "parse RSA private key failed")
	}

	p := jwt.NewJWTTokenProvider(key)
	misc.SetPublicKey(key.PublicKey)
	// Its time to initialize any module
	user.SetProvider(p)
	return nil
}
