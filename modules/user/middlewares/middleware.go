package middlewares

import (
	"context"
	"fmt"

	"github.com/fzerorubigd/balloon/pkg/assert"

	"github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/pkg/errors"

	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/resources"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type contextKey int

const (
	resource contextKey = 0
	user     contextKey = 1
)

func requirement(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, ok := resources.QueryResource(info.FullMethod)
	if ok {
		ctx = context.WithValue(ctx, resource, res)
	}

	return handler(ctx, req)
}

func auth(ctx context.Context) (context.Context, error) {
	r := ctx.Value(resource)
	if r == nil { // No user requested here
		return ctx, nil
	}
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	fmt.Println(err, "==>", token)
	fmt.Println("Key => ", ctx.Value(resource))
	return ctx, nil
}

// ExtractUser try to extract the current user from the context
func ExtractUser(ctx context.Context) (*userpb.User, error) {
	u, ok := ctx.Value(user).(*userpb.User)
	if !ok {
		return nil, errors.New("no user in context")
	}
	return u, nil
}

// MustExtractUser return the current user and panics when the user is not found
func MustExtractUser(ctx context.Context) *userpb.User {
	u, err := ExtractUser(ctx)
	assert.Nil(err)
	return u
}

func init() {
	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Unary:  requirement,
		Stream: nil, // TODO : Stream?
	})
	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Stream: grpc_auth.StreamServerInterceptor(auth),
		Unary:  grpc_auth.UnaryServerInterceptor(auth),
	})
}
