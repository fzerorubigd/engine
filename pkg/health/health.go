package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	lock sync.RWMutex
	all  = map[string]Interface{}
)

// Interface is the checker interface
type Interface interface {
	// Check must return an error if the health is not ok
	Healthy(context.Context) error
}

type healthErr map[string]error

func (h healthErr) Status() int {
	return http.StatusInternalServerError
}

func (h healthErr) Message() string {
	return "Health check failed"
}

func (h healthErr) Fields() map[string]string {
	t := make(map[string]string, len(h))
	for i := range h {
		t[i] = h[i].Error()
	}

	return t
}

func (h healthErr) Error() string {
	ret := ""
	for i := range h {
		ret += fmt.Sprintf("%s: %s\n", i, h[i])
	}

	return ret
}

// Healthy run all health checks one by one and return the errors.
func Healthy(ctx context.Context) error {
	lock.RLock()
	defer lock.RUnlock()

	var (
		m healthErr
	)
	for i := range all {
		nCtx, cnl := context.WithTimeout(ctx, time.Second)
		e := all[i].Healthy(nCtx)
		if e != nil {
			if m == nil {
				m = make(healthErr)
			}
			m[i] = e
		}
		cnl()
	}

	var err error
	if len(m) > 0 {
		err = m
	}

	return err
}

// Register a new health reader system
func Register(name string, h Interface) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := all[name]; ok {
		panic("the name is already registered")
	}

	all[name] = h
}
