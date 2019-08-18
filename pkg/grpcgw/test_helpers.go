package grpcgw

import (
	"elbix.dev/engine/pkg/log"
	"github.com/fullstorydev/grpchan/inprocgrpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

// GRPCChannel is a helper function, it is exported for tests only, do not use it!
func GRPCChannel() *inprocgrpc.Channel {
	unaryMiddle := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(log.Logger()),
	}

	streamMiddle := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(log.Logger()),
	}

	for i := range interceptors {
		if interceptors[i].Unary != nil {
			unaryMiddle = append(unaryMiddle, interceptors[i].Unary)
		}
		if interceptors[i].Stream != nil {
			streamMiddle = append(streamMiddle, interceptors[i].Stream)
		}
	}
	c := &inprocgrpc.Channel{}
	c = c.WithServerUnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMiddle...)).
		WithServerStreamInterceptor(grpc_middleware.ChainStreamServer(streamMiddle...))

	return c
}
