package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT,
}

// Context returns a context that is cancelled automatically when a kill signal received
func Context() context.Context {
	var sig = make(chan os.Signal, len(signals))
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, signals...)
	go func() {
		<-sig
		cancel()
	}()

	return ctx
}
