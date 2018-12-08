package impl

import (
	"context"

	"github.com/fzerorubigd/balloon/modules/user/proto"

	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userController struct {
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (uc *userController) Login(ctx context.Context, lr *userpb.LoginRequest) (*userpb.UserResponse, error) {
	// create and send header
	header := metadata.Pairs("Grpc-Metadata-Header-Key", "val",
		"Grpc-Trailer-Header-Key2", "val2",
		"grpcgateway-Header-Key3", "val3",
	)
	if err := grpc.SendHeader(ctx, header); err != nil {
		panic(err)
	}
	// create and set trailer
	trailer := metadata.Pairs("Grpc-Metadata-t-Key", "val",
		"Grpc-Trailer-t-Key2", "val2",
		"grpcgateway-t-Key3", "val3",
	)

	if err := grpc.SetTrailer(ctx, trailer); err != nil {
		panic(err)
	}

	return &userpb.UserResponse{
		Email:  "hi@me.com",
		Id:     "sss",
		Status: userpb.UserStatus_USER_STATUS_ACTIVE,
	}, nil
}

func (uc *userController) Logout(context.Context, *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	panic("implement me")
}

func (uc *userController) Register(context.Context, *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	panic("implement me")
}

func init() {
	grpcgw.Register(userpb.NewWrappedUserSystemServer(&userController{}))
}