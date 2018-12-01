package userpb

import (
	"context"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gopkg.in/go-playground/validator.v9"
)

type userController struct {
	v *validator.Validate
}

func (uc *userController) Initialize(ctx context.Context) {
}

func (uc *userController) Login(ctx context.Context, lr *LoginRequest) (*UserResponse, error) {
	if err := uc.v.Struct(lr); err != nil {
		return nil, err
	}

	return &UserResponse{
		Email:    "hi@me.com",
		Id:       "sss",
		Status:   UserStatus_USER_STATUS_ACTIVE,
		Username: "hi",
	}, nil
}

func (*userController) Logout(context.Context, *LogoutRequest) (*NoopResponse, error) {
	return &NoopResponse{}, nil
}

func (uc *userController) Init(ctx context.Context, ch inprocgrpc.Channel, mux *runtime.ServeMux) {
	RegisterHandlerUserSystem(&ch, uc)
	cl := NewUserSystemChannelClient(&ch)

	assert.Nil(RegisterUserSystemHandlerClient(ctx, mux, cl))
}

func init() {
	grpcgw.Register(&userController{v: validator.New()})
}
