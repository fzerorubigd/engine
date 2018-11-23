package userimpl

import (
	"context"

	"github.com/fzerorubigd/balloon/modules/user/proto"
)

type userController struct {
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (*userController) Login(context.Context, *userpb.LoginRequest) (*userpb.User, error) {
	return &userpb.User{
		Email:    "hi@me.com",
		Id:       100,
		Status:   userpb.User_USER_STATUS_ACTIVE,
		Username: "hi",
	}, nil
}

func (*userController) Logout(context.Context, *userpb.LogoutRequest) (*userpb.NoopResponse, error) {
	return &userpb.NoopResponse{}, nil
}

func NewUserController() userpb.UserSystemServer {
	return &userController{}
}
