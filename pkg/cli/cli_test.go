package cli

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCliContext(t *testing.T) {
	ctx := Context(syscall.SIGUSR1)
	select {
	case <-ctx.Done():
		require.True(t, false, "context canceled early")
	default:
	}
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		require.True(t, false, "context not canceled")
	}

	// Retry with new context
	ctx = Context(syscall.SIGUSR1)
	select {
	case <-ctx.Done():
		require.True(t, false, "context canceled early")
	default:
	}
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		require.True(t, false, "context not canceled")
	}
}
