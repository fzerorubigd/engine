package impl

import (
	"context"
	"fmt"

	"github.com/fzerorubigd/balloon/pkg/resources"

	"google.golang.org/grpc"

	"github.com/fzerorubigd/balloon/pkg/grpcgw"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type contextKey int

const resource contextKey = 0

func requirement(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, ok := resources.QueryResource(info.FullMethod)
	if ok {
		ctx = context.WithValue(ctx, resource, res)
	}

	return handler(ctx, req)
}

func auth(ctx context.Context) (context.Context, error) {
	fmt.Println("=======")
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	fmt.Println(err, "==>", token)
	fmt.Println("Key => ", ctx.Value(resource))
	return ctx, nil
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
