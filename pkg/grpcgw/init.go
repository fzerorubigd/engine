package grpcgw

import (
	"context"
	"net/http"
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
)

// Controller is the simple controller interface
type Controller interface {
	// Init the controller and register them in the server mux
	Init(context.Context, *inprocgrpc.Channel, *runtime.ServeMux)
}

type Interceptor struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

var (
	all          []Controller
	interceptors []Interceptor
	lock         sync.RWMutex

	addr = config.RegisterString("grpcw.http.addr", ":8090", "http address to listen to")
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

func Serve(ctx context.Context) {
	lock.RLock()
	defer lock.RUnlock()

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

	var (
		c         = &inprocgrpc.Channel{}
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux()
	)

	c = c.WithServerUnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMiddle...)).
		WithServerStreamInterceptor(grpc_middleware.ChainStreamServer(streamMiddle...))

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
