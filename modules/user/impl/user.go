package impl

import (
	"context"

	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/pkg/errors"
)

type userController struct {
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (uc *userController) Login(ctx context.Context, lr *userpb.LoginRequest) (*userpb.UserResponse, error) {
	m := userpb.NewManager()

	u, err := m.LoginUserByPassword(ctx, lr.GetEmail(), lr.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "email and/or password is wrong")
	}

	resp := userpb.UserResponse{
		Email:  u.GetEmail(),
		Status: u.GetStatus(),
		Id:     u.GetId(),
	}

	return &resp, nil
}

func (uc *userController) Logout(context.Context, *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	panic("S")
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
	}, nil
}

func init() {
	grpcgw.Register(userpb.NewWrappedUserSystemServer(&userController{}))
}
