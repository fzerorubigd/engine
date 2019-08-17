package sentry

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/initializer"
	"elbix.dev/engine/pkg/log"
	"elbix.dev/engine/pkg/version"
)

var (
	sentryAddress = config.RegisterString(
		"services.sentry.url",
		"https://sentry.elbix.dev",
		"Sentry domain",
	)

	sentryProject = config.RegisterString(
		"services.sentry.project",
		"2",
		"Sentry project id",
	)

	sentryProjectKey = config.RegisterString(
		"services.sentry.secret",
		"",
		"Sentry secret",
	)

	sentryEnabled = config.RegisterBoolean(
		"services.sentry.enabled",
		false,
		"Sentry enabled or not",
	)

	enabled bool
)

type sentryInit struct {
}

func (s *sentryInit) Initialize(ctx context.Context) {
	if !sentryEnabled.Bool() {
		enabled = false
		return
	}
	dsn, err := url.Parse(sentryAddress.String())
	if err != nil {
		log.Info("Sentry URL is invalid", log.Err(err))
		enabled = false
		return
	}

	dsn.User = url.User(sentryProjectKey.String())
	dsn.Path = "/" + sentryProject.String()

	if err := sentrygo.Init(sentrygo.ClientOptions{
		Dsn:     dsn.String(),
		Release: fmt.Sprintf("BUILD-%d", version.GetVersion().Count),
	}); err != nil {
		enabled = false
		log.Info("Sentry initialization failed", log.Err(err))
		return
	}
	enabled = true
	log.Info("Sentry integration enabled")
	go func() {
		<-ctx.Done()
		sentrygo.Flush(time.Second)
	}()
}

type withStack struct {
	p   interface{}
	pcs []uintptr
}

func (w *withStack) StackTrace() []uintptr {
	return w.pcs
}

func (w *withStack) Error() string {
	switch t := w.p.(type) {
	case error:
		return t.Error()
	default:
		return fmt.Sprint(t)
	}
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(5, pcs[:])
	return pcs[0:n]
}

// Recover convert the recovered p to an error and also send it to sentry
func Recover(p interface{}) error {
	log.Error("Recover from panic", log.Any("panic", p))
	if enabled {
		w := &withStack{
			p:   p,
			pcs: callers(),
		}
		go sentrygo.CaptureException(w)
	}
	return status.Errorf(codes.Internal, "%s", p)
}

func init() {
	initializer.Register(&sentryInit{}, 0)
}
