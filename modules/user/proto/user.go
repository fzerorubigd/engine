package userpb

import (
	"context"
	"errors"
)

type userController struct {
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (*userController) Login(ctx context.Context, lr *LoginRequest) (*User, error) {
	if lr == nil || lr.Username != "aaa" {
		return nil, errors.New("Sss")
	}
	return &User{
		Email:    "hi@me.com",
		Id:       100,
		Status:   User_USER_STATUS_ACTIVE,
		Username: "hi",
	}, nil
}

func (*userController) Logout(context.Context, *LogoutRequest) (*NoopResponse, error) {
	return &NoopResponse{}, nil
}

func NewUserController() UserSystemServer {
	return &userController{}
}
