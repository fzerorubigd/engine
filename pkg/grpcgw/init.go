package grpcgw

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/fzerorubigd/balloon/pkg/config"
	"github.com/fzerorubigd/balloon/pkg/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
)

// Controller is the simple controller interface
type Controller interface {
	// Init the controller and register them in the server mux
	Init(context.Context, inprocgrpc.Channel, *runtime.ServeMux)
}

var (
	all  []Controller
	lock sync.RWMutex

	addr = config.RegisterString("grpcw.http.addr", ":8000", "http address to listen to")
)

// Register new controller into system
func Register(c Controller) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, c)
}

func Serve(ctx context.Context) {
	lock.RLock()
	defer lock.RUnlock()

	var (
		c         inprocgrpc.Channel
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux()
	)

	normalMux.HandleFunc("/swagger", swaggerHandler)
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
