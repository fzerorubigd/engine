package impl

import (
	"context"
	"time"

	"elbix.dev/engine/modules/user/middlewares"
	"elbix.dev/engine/modules/user/proto"
	"elbix.dev/engine/pkg/assert"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/log"
	"elbix.dev/engine/pkg/token"
)

var (
	expire   = config.RegisterDuration("modules.user.token.expire", time.Hour*24*3, "token expiration timeout")
	provider token.Provider
)

type userController struct {
}

func (uc *userController) ForgotPassword(ctx context.Context, fp *userpb.ForgotPasswordRequest) (*userpb.ForgotPasswordResponse, error) {
	m := userpb.NewManager()

	// TODO : rate limit
	u, err := m.FindUserByEmail(ctx, fp.Email)
	// This is a little tricky here, even on error we return an ok, the message on the client
	// should be complete
	if err != nil {
		return &userpb.ForgotPasswordResponse{}, nil
	}

	token, _, err := m.CreateForgottenToken(ctx, u)
	assert.Nil(err)

	log.Info("TODO", log.String("token", token))

	return &userpb.ForgotPasswordResponse{}, nil
}

func (uc *userController) ChangeDisplayName(ctx context.Context, cd *userpb.ChangeDisplayNameRequest) (*userpb.ChangeDisplayNameResponse, error) {
	u := middlewares.MustExtractUser(ctx)
	m := userpb.NewManager()
	u, err := m.GetUserByPrimary(ctx, u.Id)
	assert.Nil(err)
	u.DisplayName = cd.DisplayName
	assert.Nil(m.UpdateUser(ctx, u))
	m.UpdateToken(ctx, u, expire.Duration(), middlewares.MustExtractToken(ctx))
	return &userpb.ChangeDisplayNameResponse{}, nil
}

func (uc *userController) ChangePassword(ctx context.Context, cpr *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	old := middlewares.MustExtractUser(ctx)
	// ok reload the user from db
	m := userpb.NewManager()
	u, err := m.FindUserByEmailPassword(ctx, old.GetEmail(), cpr.GetOldPassword())
	if err != nil {
		return nil, grpcgw.NewBadRequest(err, "old password is wrong")
	}

	assert.Nil(m.ChangePassword(ctx, u, cpr.GetNewPassword()))

	return &userpb.ChangePasswordResponse{}, nil
}

func (uc *userController) Ping(ctx context.Context, _ *userpb.PingRequest) (*userpb.UserResponse, error) {
	u := middlewares.MustExtractUser(ctx)
	tok := middlewares.MustExtractToken(ctx)

	return &userpb.UserResponse{
		Id:          u.GetId(),
		Token:       tok,
		Status:      u.GetStatus(),
		DisplayName: u.DisplayName,
	}, nil
}

func (uc *userController) Login(ctx context.Context, lr *userpb.LoginRequest) (*userpb.UserResponse, error) {
	m := userpb.NewManager()

	u, err := m.FindUserByEmailPassword(ctx, lr.GetEmail(), lr.GetPassword())
	if err != nil {
		return nil, grpcgw.NewBadRequest(err, "email and/or password is wrong")
	}

	resp := userpb.UserResponse{
		DisplayName: u.GetDisplayName(),
		Status:      u.GetStatus(),
		Id:          u.GetId(),
		Token:       m.CreateToken(ctx, u, expire.Duration()),
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

	u, err := m.RegisterUser(ctx, ru.GetEmail(), ru.GetDisplayName(), ru.GetPassword())
	if err != nil {
		return nil, grpcgw.NewBadRequest(err, "duplicate email")
	}

	return &userpb.UserResponse{
		Id:          u.GetId(),
		Status:      u.GetStatus(),
		DisplayName: u.GetDisplayName(),
		Token:       m.CreateToken(ctx, u, expire.Duration()),
	}, nil
}

func (uc *userController) Initialize(ctx context.Context) {
}

// SetProvider for setting the token provider
func SetProvider(p token.Provider) {
	provider = p
}

func init() {
	grpcgw.Register(userpb.NewWrappedUserSystemServer(&userController{}))
}
