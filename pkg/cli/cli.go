package cli

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
)

var (
	lock = sync.Mutex{}
	old  = map[string]context.Context{}
)

// Context returns a context that is cancelled automatically when a kill
// signal received
// The idea is to create one context based on a combination of signalls, and
// not one context for each call.
func Context(signals ...os.Signal) context.Context {
	lock.Lock()
	defer lock.Unlock()

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

	key := signalKey(signals)
	if ctx, ok := old[key]; ok {
		return ctx
	}

	var sig = make(chan os.Signal, len(signals))
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, signals...)
	go func() {
		<-sig
		cancel()

		// Make sure we delete the key, so the next time -IF- someone create
		// the same watch, they get a new context this time (not a canceled context)
		lock.Lock()
		defer lock.Unlock()

		delete(old, key)
	}()

	old[key] = ctx
	return ctx
}

func signalKey(signals []os.Signal) string {
	sort.Slice(signals, func(i, j int) bool {
		return strings.Compare(signals[i].String(), signals[j].String()) < 0
	})

	var key string
	for i := range signals {
		key += signals[i].String()
	}

	return key
}
