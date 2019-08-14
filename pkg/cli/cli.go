package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var sig = make(chan os.Signal, 4)

// Context returns a context that is cancelled automatically when a kill signal received
func Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)
	go func() {
		<-sig
		cancel()
	}()

	return ctx
}
