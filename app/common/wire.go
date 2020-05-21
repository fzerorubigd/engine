package common

/*
This file is for wiring the different modules in the project.
*/

import (
	miscimpl "elbix.dev/engine/modules/misc/impl"
	miscpb "elbix.dev/engine/modules/misc/proto"
	usrimpl "elbix.dev/engine/modules/user/impl"
	userpb "elbix.dev/engine/modules/user/proto"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/sec"
	"elbix.dev/engine/pkg/token/jwt"
)

func userMod() (grpcgw.Controller, error) {
	privateKey := getPrivateKey()
	rsaPrivateKey, err := sec.ParseRSAPrivateKeyFromBase64PEM(privateKey)
	if err != nil {
		return nil, err
	}
	provider := jwt.NewJWTTokenProvider(rsaPrivateKey)
	userSystemServer := usrimpl.NewUserController(provider)
	wrappedUserSystemController := userpb.NewWrappedUserSystemServer(userSystemServer)
	return wrappedUserSystemController, nil
}

func miscMod() (grpcgw.Controller, error) {
	privateKey := getPrivateKey()
	rsaPrivateKey, err := sec.ParseRSAPrivateKeyFromBase64PEM(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := sec.ExtractPublicFromPrivate(rsaPrivateKey)
	miscSystemServer := miscimpl.NewMiscController(publicKey)
	wrappedMiscSystemController := miscpb.NewWrappedMiscSystemServer(miscSystemServer)
	return wrappedMiscSystemController, nil
}
