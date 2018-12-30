package grpcgw

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/pkg/config"
	"github.com/fzerorubigd/balloon/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"gopkg.in/fzerorubigd/onion.v3"
)

// Controller is the simple controller interface
type Controller interface {
	// Init the controller and register them in the server mux
	Init(context.Context, *inprocgrpc.Channel, *runtime.ServeMux)
}

// Interceptor type used to register both interceptors at the same time
type Interceptor struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

var (
	all          []Controller
	interceptors []Interceptor
	lock         sync.RWMutex

	addr onion.String
)

// Register new controller into system
func Register(c Controller) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, c)
}

// RegisterInterceptors register middleware
func RegisterInterceptors(i Interceptor) {
	interceptors = append(interceptors, i)
}

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

// Serve start the server and wait
func Serve(ctx context.Context) {
	lock.RLock()
	defer lock.RUnlock()

	var (
		c         = GRPCChannel()
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux()
	)

	normalMux.HandleFunc("/v1/swagger/", swaggerHandler)
	for i := range all {
		all[i].Init(ctx, c, mux)
	}

	// TODO : Handle cors with options, currently it is ok.
	normalMux.Handle("/", cors.AllowAll().Handler(mux))
	srv := http.Server{
		Addr:    addr.String(),
		Handler: normalMux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("Http listen and serve failed", log.Err(err))
		}
	}()

	<-ctx.Done()
	nCtx, cnl := context.WithTimeout(context.Background(), time.Second)
	defer cnl()

	if err := srv.Shutdown(nCtx); err != nil {
		log.Error("Server shutdown failed", log.Err(err))
	}
}

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	addr = config.RegisterString(
		"grpcw.http.addr",
		fmt.Sprintf(":%s", port),
		"http address to listen to",
	)
}
