package health

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"elbix.dev/engine/pkg/grpcgw"
)

type healthy struct {
	err error
}

func (h *healthy) Healthy(ctx context.Context) error {
	return h.err
}

func TestRegister(t *testing.T) {
	errs := Healthy(context.Background())
	require.NoError(t, errs)
	h1 := &healthy{}
	require.NotPanics(t, func() { Register("myname", h1) })
	require.Panics(t, func() { Register("myname", h1) })

	errs = Healthy(context.Background())
	require.NoError(t, errs)
	h1.err = errors.New("err")
	errs = Healthy(context.Background())
	assert.Error(t, errs)

	expected := healthErr{
		"myname": h1.err,
	}
	assert.Equal(t, expected, errs)
	assert.Implements(t, (*grpcgw.GWError)(nil), errs)
	assert.Equal(t, http.StatusInternalServerError, errs.(grpcgw.GWError).Status())
	assert.NotEmpty(t, errs.(grpcgw.GWError).Message())
	assert.Equal(t, map[string]string{
		"myname": "err",
	}, errs.(grpcgw.GWError).Fields())
	assert.Equal(t, "myname: err\n", errs.Error())
}
