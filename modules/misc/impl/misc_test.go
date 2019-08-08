package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	miscpb "github.com/fzerorubigd/engine/modules/misc/proto"
	"github.com/fzerorubigd/engine/pkg/grpcgw"
	"github.com/fzerorubigd/engine/pkg/health"
	"github.com/fzerorubigd/engine/pkg/version"
)

var ch *inprocgrpc.Channel

func newClient() miscpb.MiscSystemClient {
	if ch == nil {
		ch = grpcgw.GRPCChannel()
		miscpb.RegisterHandlerMiscSystem(ch, miscpb.NewWrappedMiscSystemServer(&miscController{}))
	}
	return miscpb.NewMiscSystemChannelClient(ch)
}

type healthFake struct {
	err error
}

func (h *healthFake) Healthy(context.Context) error {
	return h.err
}

func TestMiscController_Health(t *testing.T) {
	ctx := context.Background()
	// defer mockery.Start(ctx, t)()
	h := &healthFake{}
	health.Register("fake", h)

	m := newClient()
	r, err := m.Health(ctx, &miscpb.HealthRequest{})
	assert.NotNil(t, r)
	assert.NoError(t, err)

	h.err = errors.New("err")

	r, err = m.Health(ctx, &miscpb.HealthRequest{})
	assert.Nil(t, r)
	assert.Error(t, err)
	require.Implements(t, (*grpcgw.GWError)(nil), err)
}

func TestMiscController_Version(t *testing.T) {
	ctx := context.Background()
	// defer mockery.Start(ctx, t)()
	m := newClient()
	v, err := m.Version(ctx, &miscpb.VersionRequest{})
	require.NoError(t, err)
	require.NotNil(t, v)

	v2 := version.GetVersion()
	require.Equal(t, v2.Count, v.Count)
	require.Equal(t, v2.Hash, v.CommitHash)
	require.Equal(t, v2.Short, v.ShortHash)
	require.Equal(t, v2.BuildDate.Unix(), v.BuildDate.GetSeconds())
	require.Equal(t, v2.Date.Unix(), v.CommitDate.GetSeconds())
}
