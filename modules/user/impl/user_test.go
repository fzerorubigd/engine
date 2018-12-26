package impl

import (
	"context"
	"net/http"
	"testing"

	"github.com/fullstorydev/grpchan/inprocgrpc"

	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/fzerorubigd/balloon/pkg/mockery"
)

var ch *inprocgrpc.Channel

func newClient() userpb.UserSystemClient {
	if ch == nil {
		ch = grpcgw.GRPCChannel()
		userpb.RegisterHandlerUserSystem(ch, userpb.NewWrappedUserSystemServer(&userController{}))
	}
	return userpb.NewUserSystemChannelClient(ch)
}

func TestUserController_Login_Invalid(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r, err := u.Login(ctx, &userpb.LoginRequest{
		Email:    "invalid_email",
		Password: "123",
	})
	assert.Nil(t, r)
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr := err.(grpcgw.GWError)
	assert.Equal(t, map[string]string{"Email": "email", "Password": "gte"}, gErr.Fields())
	assert.Equal(t, http.StatusBadRequest, gErr.Status())

	r, err = u.Login(ctx, &userpb.LoginRequest{
		Email:    "valid@email.com",
		Password: "validpass",
	})
	assert.Nil(t, r)
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr = err.(grpcgw.GWError)
	require.Nil(t, gErr.Fields())
	assert.Equal(t, http.StatusBadRequest, gErr.Status())
}

func TestUserController_Register_Invalid(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r, err := u.Register(ctx, &userpb.RegisterRequest{
		Email:    "invalid_email",
		Password: "123",
	})
	assert.Nil(t, r)
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr := err.(grpcgw.GWError)
	assert.Equal(t, map[string]string{"Email": "email", "Password": "gte"}, gErr.Fields())
	assert.Equal(t, http.StatusBadRequest, gErr.Status())

	r, err = u.Register(ctx, &userpb.RegisterRequest{
		Email:    "master@cerulean.ir", // Email from migration
		Password: "12345678",
	})
	assert.Nil(t, r)
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr = err.(grpcgw.GWError)
	assert.Nil(t, gErr.Fields())
	assert.Equal(t, http.StatusBadRequest, gErr.Status())

}

func TestUserController_Register(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r, err := u.Register(ctx, &userpb.RegisterRequest{
		Email:    "valid@gmail.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.Equal(t, "valid@gmail.com", r.Email)
	assert.Equal(t, userpb.UserStatus_USER_STATUS_REGISTERED, r.Status)
	assert.NotZerof(t, r.Id, "User Id is not valid")

	r2, err := u.Login(ctx, &userpb.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.Equal(t, r.Id, r2.Id)
	assert.Equal(t, r.Status, r2.Status)
	assert.Equal(t, r.Email, r2.Email)
}

func TestUserController_Logout(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r, err := u.Logout(ctx, &userpb.LogoutRequest{})
	assert.Nil(t, r)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr := err.(grpcgw.GWError)
	assert.Equal(t, http.StatusUnauthorized, gErr.Status())

	r1, err := u.Register(ctx, &userpb.RegisterRequest{
		Email:    "valid@gmail.com",
		Password: "bita123",
	})
	assert.NoError(t, err)

	_, err = u.Logout(mockery.AuthorizeToken(ctx, r1.Token), &userpb.LogoutRequest{})
	assert.NoError(t, err)

	_, err = u.Logout(mockery.AuthorizeToken(ctx, r1.Token), &userpb.LogoutRequest{})
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr = err.(grpcgw.GWError)
	assert.Equal(t, http.StatusUnauthorized, gErr.Status())
}

func TestUserController_Ping(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r, err := u.Logout(ctx, &userpb.LogoutRequest{})
	assert.Nil(t, r)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr := err.(grpcgw.GWError)
	assert.Equal(t, http.StatusUnauthorized, gErr.Status())

	r1, err := u.Register(ctx, &userpb.RegisterRequest{
		Email:    "valid@gmail.com",
		Password: "bita123",
	})
	assert.NoError(t, err)

	r2, err := u.Ping(mockery.AuthorizeToken(ctx, r1.Token), &userpb.PingRequest{})
	require.NoError(t, err)
	assert.Equal(t, r1, r2)
}

func TestUserController_ChangePassword(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(t, ctx)()

	u := newClient()
	r1, err := u.Register(ctx, &userpb.RegisterRequest{
		Email:    "valid@gmail.com",
		Password: "bita123",
	})
	assert.NoError(t, err)

	ctx = mockery.AuthorizeToken(ctx, r1.Token)

	r2, err := u.ChangePassword(ctx, &userpb.ChangePasswordRequest{
		OldPassword: "wrongpass",
		NewPassword: "newpass",
	})

	assert.Nil(t, r2)
	assert.Error(t, err)
	require.IsType(t, grpcgw.NewNotFound(nil), err)
	gErr := err.(grpcgw.GWError)
	assert.Equal(t, http.StatusBadRequest, gErr.Status())

	r2, err = u.ChangePassword(ctx, &userpb.ChangePasswordRequest{
		OldPassword: "bita123",
		NewPassword: "newpass",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r2)

	r3, err := u.Login(ctx, &userpb.LoginRequest{
		Password: "newpass",
		Email:    "valid@gmail.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, r3.Id, r1.Id)
	assert.Equal(t, r3.Email, r1.Email)
	assert.Equal(t, r3.Status, r1.Status)
}
