package cli

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCliContext(t *testing.T) {
	signals = append(signals, syscall.SIGUSR1)
	ctx := Context()
	select {
	case <-ctx.Done():
		require.True(t, false, "context canceled early")
	default:
	}
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		require.True(t, false, "context not canceled")
	}
}
