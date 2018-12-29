package cli

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	ctx := Context()
	select {
	case <-ctx.Done():
		require.True(t, false, "context canceled early")
	default:
	}
	sig <- syscall.SIGINT

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		require.True(t, false, "context not canceled")
	}

}
