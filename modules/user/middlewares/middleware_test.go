package middlewares

import (
	"context"
	"testing"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/mockery"
	"github.com/stretchr/testify/assert"
)

type userMock struct {
	t     *testing.T
	token string
	u     *userpb.User
}

func (um userMock) ChangeDisplayName(context.Context, *userpb.ChangeDisplayNameRequest) (*userpb.ChangeDisplayNameResponse, error) {
	panic("implement me")
}

func (um userMock) Login(ctx context.Context, _ *userpb.LoginRequest) (*userpb.UserResponse, error) {
	// this is not protected route
	u, err := ExtractUser(ctx)
	assert.Nil(um.t, u)
	assert.Error(um.t, err)

	t, err := ExtractToken(ctx)
	assert.Empty(um.t, t)
	assert.Error(um.t, err)

	return &userpb.UserResponse{}, nil
}

func (userMock) Logout(context.Context, *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	panic("implement me")
}

func (userMock) Register(context.Context, *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	panic("implement me")
}

func (um userMock) Ping(ctx context.Context, pr *userpb.PingRequest) (*userpb.UserResponse, error) {
	u := MustExtractUser(ctx)
	// Do not check equality of entire object because of the grpc extra fields
	assert.Equal(um.t, um.u.Id, u.Id)
	assert.Equal(um.t, um.u.Email, u.Email)
	t := MustExtractToken(ctx)
	assert.Equal(um.t, um.token, t)

	// All fields are not required
	return &userpb.UserResponse{
		DisplayName: um.u.DisplayName,
		Id:          um.u.Id,
		Token:       um.token,
	}, nil
}

func (userMock) ChangePassword(context.Context, *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	panic("implement me")
}

var ch *inprocgrpc.Channel

func newClient(u *userMock) userpb.UserSystemClient {
	if ch == nil {
		ch = grpcgw.GRPCChannel()
		userpb.RegisterHandlerUserSystem(ch, userpb.NewWrappedUserSystemServer(u))
	}
	return userpb.NewUserSystemChannelClient(ch)
}

func TestMiddlewareSystem(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(ctx, t)()

	user := &userpb.User{
		Email:       "valid@email.com",
		DisplayName: "display",
		Id:          1,
	}

	mock := &userMock{
		token: userpb.NewManager().CreateToken(ctx, user, time.Hour),
		t:     t,
		u:     user,
	}

	cl := newClient(mock)
	nctx := mockery.AuthorizeToken(ctx, mock.token)

	r, err := cl.Ping(nctx, &userpb.PingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, r.Id, user.Id)
	assert.Equal(t, r.DisplayName, user.DisplayName)
	assert.Equal(t, r.Token, mock.token)

	r1, err := cl.Login(ctx, &userpb.LoginRequest{
		Email:    user.Email,
		Password: "123456", // Pass the validation
	})
	assert.NoError(t, err)
	assert.NotNil(t, r1)

	// No token
	r, err = cl.Ping(ctx, &userpb.PingRequest{})
	assert.Nil(t, r)
	assert.Error(t, err)

	// Invalid token
	ictx := mockery.AuthorizeToken(ctx, "INVALID_TOKEN")
	r, err = cl.Ping(ictx, &userpb.PingRequest{})
	assert.Nil(t, r)
	assert.Error(t, err)

}
