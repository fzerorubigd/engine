package grpcgw

import "context"

// Middleware is the middleware system for the grpc-gateway
type Middleware interface {
	Transform(ctx context.Context) (context.Context, error)
}

// ExecuteMiddleware try to execute middleware on an interface
func ExecuteMiddleware(ctx context.Context, in interface{}) (context.Context, error) {
	if m, ok := in.(Middleware); ok {
		return m.Transform(ctx)
	}

	return ctx, nil
}
