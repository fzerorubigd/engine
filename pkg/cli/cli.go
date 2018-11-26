package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// sig is a side channel that may receive signals when ReceiveSignal. Any
// context that was created via cli.Context will wait for signals on this
// channel.
var sig = make(chan os.Signal, 4)

// Context returns a context that is cancelled automatically when a SIGINT,
// SIGQUIT or SIGTERM signal is received.
func Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)
	go func() {
		select {
		case <-sig:
			cancel()
		}
	}()

	return ctx
}

// ReceiveSignal propagates the given signal to all contexts that may have been
// created via cli.Context(). This function is only useful when the terminal is
// hijacked and you want to emulate signals manually.
func ReceiveSignal(s os.Signal) {
	select {
	case sig <- s:
	default:
	}
}
