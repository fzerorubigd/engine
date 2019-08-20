package impl

import (
	"context"
	"errors"
	"time"

	"elbix.dev/engine/modules/user/middlewares"
	userpb "elbix.dev/engine/modules/user/proto"
	"elbix.dev/engine/pkg/assert"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/grpcgw"
	"elbix.dev/engine/pkg/log"
	"elbix.dev/engine/pkg/token"
	"google.golang.org/api/oauth2/v2"
	"gopkg.in/go-playground/validator.v9"
)

var (
	expire = config.RegisterDuration("modules.user.token.expire", time.Hour*24*3, "token expiration timeout")
)

type userController struct {
	v *validator.Validate
	p token.Provider
}

func (uc *userController) VerifyToken(ctx context.Context, vt *userpb.VerifyTokenRequest) (*userpb.UserResponse, error) {
	oauth2Service, err := oauth2.NewService(ctx)
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(vt.TokenId)
	tok, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	assert.Nil(uc.v.VarCtx(ctx, tok.Email, "required,email"))

	m := userpb.NewManager()

	u, err := m.FindUserByEmail(ctx, tok.Email)
	if err != nil {
		u, err = m.RegisterUser(ctx, tok.Email, "New User", userpb.NoPassString)
		if err != nil {
			return nil, grpcgw.NewBadRequest(err, "something is wrong")
		}
	}

	return &userpb.UserResponse{
		Id:             u.GetId(),
		Status:         u.GetStatus(),
		DisplayName:    u.GetDisplayName(),
		Token:          m.CreateToken(ctx, u, expire.Duration()),
		ChangePassword: u.ShouldChangePass(),
	}, nil
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

	return &userpb.ChangeDisplayNameResponse{}, nil
}

func (uc *userController) ChangePassword(ctx context.Context, cpr *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	u := middlewares.MustExtractUser(ctx)
	// ok reload the user from db
	m := userpb.NewManager()
	if u.GetPassword() != userpb.NoPassString {
		if !u.VerifyPassword(cpr.GetOldPassword()) {
			return nil, grpcgw.NewBadRequest(errors.New("old pass is wrong"), "old password is wrong")
		}
	}
	assert.Nil(m.ChangePassword(ctx, u, cpr.GetNewPassword()))

	return &userpb.ChangePasswordResponse{}, nil
}

func (uc *userController) Ping(ctx context.Context, _ *userpb.PingRequest) (*userpb.UserResponse, error) {
	u1 := middlewares.MustExtractUser(ctx)
	tok := middlewares.MustExtractToken(ctx)
	m := userpb.NewManager()
	u, err := m.GetUserByPrimary(ctx, u1.GetId())
	assert.Nil(err)
	return &userpb.UserResponse{
		Id:             u.GetId(),
		Token:          tok,
		Status:         u.GetStatus(),
		DisplayName:    u.DisplayName,
		ChangePassword: u.ShouldChangePass(),
	}, nil
}

func (uc *userController) Login(ctx context.Context, lr *userpb.LoginRequest) (*userpb.UserResponse, error) {
	m := userpb.NewManager()

	u, err := m.FindUserByEmailPassword(ctx, lr.GetEmail(), lr.GetPassword())
	if err != nil {
		return nil, grpcgw.NewBadRequest(err, "email and/or password is wrong")
	}

	resp := userpb.UserResponse{
		DisplayName:    u.GetDisplayName(),
		Status:         u.GetStatus(),
		Id:             u.GetId(),
		Token:          m.CreateToken(ctx, u, expire.Duration()),
		ChangePassword: u.ShouldChangePass(),
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
		Id:             u.GetId(),
		Status:         u.GetStatus(),
		DisplayName:    u.GetDisplayName(),
		Token:          m.CreateToken(ctx, u, expire.Duration()),
		ChangePassword: u.ShouldChangePass(),
	}, nil
}

func (uc *userController) Initialize(ctx context.Context) {
}

// NewUserController return a grpc user controller
func NewUserController(p token.Provider) userpb.UserSystemServer {
	// TODO: remove this, and use the same object
	userpb.SetProvider(p)
	return &userController{
		v: validator.New(),
		p: p,
	}
}
