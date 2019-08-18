package grpcgw

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/goraz/onion/configwatch"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/log"
	"elbix.dev/engine/pkg/sentry"
)

// Controller is the simple controller interface
type Controller interface {
	// Init the controller and register them in the server mux
	Init(context.Context, *grpc.ClientConn, *runtime.ServeMux)

	InitGRPC(context.Context, *grpc.Server)
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

	httpAddr configwatch.String
	grpcAddr configwatch.String
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

// gRPCServer is a helper function, it is exported for tests only, do not use it!
func gRPCServer() *grpc.Server {
	unaryMiddle := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(sentry.Recover)),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(log.Logger()),
	}

	streamMiddle := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(sentry.Recover)),
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
	c := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMiddle...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamMiddle...)),
	)

	return c
}

// gRPCClient creates a new GRPC client conn
func gRPCClient() (*grpc.ClientConn, error) {
	return grpc.Dial(grpcAddr.String(), grpc.WithInsecure())
}

// Serve start the server and wait
func serveHTTP(ctx context.Context) (func() error, error) {
	var (
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux()
	)
	c, err := gRPCClient()
	if err != nil {
		return nil, err
	}

	normalMux.HandleFunc("/v1/swagger/", swaggerHandler)
	for i := range all {
		all[i].Init(ctx, c, mux)
	}

	// TODO : Handle cors with options, currently it is ok.
	normalMux.Handle("/", cors.AllowAll().Handler(mux))
	srv := http.Server{
		Addr:    httpAddr.String(),
		Handler: normalMux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() error {
		nCtx, cnl := context.WithTimeout(context.Background(), time.Second)
		defer cnl()

		return srv.Shutdown(nCtx)
	}, nil
}

// Serve start the server and wait
func serveGRPC(ctx context.Context) (func() error, error) {
	srv := gRPCServer()
	for i := range all {
		all[i].InitGRPC(ctx, srv)
	}

	lis, err := net.Listen("tcp", grpcAddr.String())
	if err != nil {
		return nil, err
	}
	go func() {
		err := srv.Serve(lis)
		if err != nil {
			log.Error("Connection Closed", log.Err(err))
		}
	}()

	return lis.Close, nil
}

func Serve(ctx context.Context) error {
	lock.RLock()
	defer lock.RUnlock()

	grpcFn, err := serveGRPC(ctx)
	if err != nil {
		return err
	}
	httpFn, err := serveHTTP(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	e1 := httpFn()
	e2 := grpcFn()

	if e1 != nil {
		return e1
	}

	return e2
}

func init() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8090"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	httpAddr = config.RegisterString(
		"grpcw.http.addr",
		fmt.Sprintf(":%s", httpPort),
		"http address to listen to",
	)

	grpcAddr = config.RegisterString(
		"grpcw.grpc.addr",
		fmt.Sprintf(":%s", grpcPort),
		"gRPC address to listen to",
	)

}
