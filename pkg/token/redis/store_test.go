package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"elbix.dev/engine/pkg/kv"
	"elbix.dev/engine/pkg/random"

	"elbix.dev/engine/pkg/mockery"
	"github.com/stretchr/testify/require"
)

type wrongJSON int

func (w wrongJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("invalid")
}

func TestNewRedisTokenProvider(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(ctx, t)()

	data := map[string]interface{}{
		"str": "ABCD",
		"int": 100,
	}

	s := NewRedisTokenProvider()
	tok, err := s.Store(data, time.Hour)
	require.NoError(t, err)
	ret, err := s.Fetch(tok)
	require.NoError(t, err)
	// It's tricky. int translated as float64
	require.Equal(t, data["str"].(string), ret["str"].(string))
	require.Equal(t, data["int"], int(ret["int"].(float64)))

	ret, err = s.Fetch("INVALID_TOKEN")
	require.Error(t, err)
	require.Nil(t, ret)

	id := <-random.ID
	require.NoError(t, kv.StoreKey(id, "INVALID", time.Hour))
	ret, err = s.Fetch(id)
	require.Error(t, err)
	require.Nil(t, ret)

	var d = make(map[string]interface{})
	_, err = s.Store(d, -1)
	require.Error(t, err)

	d["w"] = wrongJSON(10)
	_, err = s.Store(d, time.Hour)
	require.Error(t, err)
}
