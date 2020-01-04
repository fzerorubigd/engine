package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context returns a context that is cancelled automatically when a kill signal received
func Context(signals ...os.Signal) context.Context {
	// I prefer to listen on specific signalls and not ALL signals.
	// in notify, passing no signal means all signals.
	if len(signals) == 0 {
		signals = []os.Signal{
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGABRT,
		}
	}
	var sig = make(chan os.Signal, len(signals))
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, signals...)
	go func() {
		<-sig
		cancel()
	}()

	return ctx
}
