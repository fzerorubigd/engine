package impl

import (
	"context"
	"time"

	"github.com/fzerorubigd/balloon/modules/user/middlewares"

	"github.com/fzerorubigd/balloon/pkg/config"

	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/pkg/errors"
)

var (
	expire = config.RegisterDuration("modules.user.token.expire", time.Hour*24*3, "token expiration timeout")
)

type userController struct {
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (uc *userController) Login(ctx context.Context, lr *userpb.LoginRequest) (*userpb.UserResponse, error) {
	m := userpb.NewManager()

	u, err := m.FindUserByEmailPassword(ctx, lr.GetEmail(), lr.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "email and/or password is wrong")
	}

	resp := userpb.UserResponse{
		Email:  u.GetEmail(),
		Status: u.GetStatus(),
		Id:     u.GetId(),
		Token:  m.CreateToken(ctx, u, expire.Duration()),
	}

	return &resp, nil
}

func (uc *userController) Logout(ctx context.Context, _ *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	tok := middlewares.MustExtractToken(ctx)
	userpb.NewManager().DeleteToken(ctx, tok)
	return &userpb.LogoutResponse{}, nil
}

func (uc *userController) Register(ctx context.Context, ru *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	m := userpb.NewManager()

	u, err := m.RegisterUser(ctx, ru.GetEmail(), ru.GetPassword())
	if err != nil {
		return nil, grpcgw.NewBadRequest(err, "duplicate email")
	}

	return &userpb.UserResponse{
		Id:     u.GetId(),
		Status: u.GetStatus(),
		Email:  u.GetEmail(),
		Token:  m.CreateToken(ctx, u, expire.Duration()),
	}, nil
}

func init() {
	grpcgw.Register(userpb.NewWrappedUserSystemServer(&userController{}))
}
