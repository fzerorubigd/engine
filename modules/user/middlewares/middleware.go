package middlewares

import (
	"context"
	"net/http"

	userpb "github.com/fzerorubigd/balloon/modules/user/proto"
	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/fzerorubigd/balloon/pkg/resources"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type contextKey int

const (
	resource contextKey = 0
	user     contextKey = 1
	token    contextKey = 2
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
	tok, err := grpc_auth.AuthFromMD(ctx, "bearer")
	m := userpb.NewManager()
	u, err := m.FindUserByIndirectToken(ctx, tok)
	if err != nil {
		return ctx, grpcgw.NewBadRequestStatus(err, "invalid token", http.StatusUnauthorized)
	}

	return context.WithValue(context.WithValue(ctx, user, u), token, tok), nil
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

func ExtractToken(ctx context.Context) (string, error) {
	tok, ok := ctx.Value(token).(string)
	if !ok {
		return "", errors.New("no token in context")
	}
	return tok, nil
}

func MustExtractToken(ctx context.Context) string {
	tok, err := ExtractToken(ctx)
	assert.Nil(err)
	return tok
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
