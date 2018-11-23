package initializer

import (
	"context"
	"sort"
	"time"
)

// Interface is the type to call early on system initialize call
type Interface interface {
	Initialize(context.Context)
}

type single struct {
	in    Interface
	order int
}

type groups []single

func (g groups) Len() int {
	return len(g)
}

func (g groups) Less(i, j int) bool {
	return g[i].order < g[j].order
}

func (g groups) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

var (
	gr = make(groups, 0)
)

// Register a module in initializer
func Register(initializer Interface, order int) {
	gr = append(gr, single{in: initializer, order: order})
}

// Initialize all modules and return the finalizer function
func Initialize(ctx context.Context) func() {
	ctx, cnl := context.WithCancel(ctx)
	sort.Sort(gr)
	for i := range gr {
		gr[i].in.Initialize(ctx)
	}

	return func() {
		cnl()
		<-time.After(time.Second)
	}
}
