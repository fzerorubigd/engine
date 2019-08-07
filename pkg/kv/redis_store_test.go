package kv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fzerorubigd/engine/pkg/mockery"
	"github.com/fzerorubigd/engine/pkg/random"
)

func TestStoreKey(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(ctx, t)()

	key1, key2 := <-random.ID, <-random.ID
	assert.NoError(t, StoreKey(key1, "value1", time.Minute))
	assert.NotPanics(t, func() { MustStoreKey(key2, "value2", time.Minute) })

	s1, err := FetchKey(key1)
	assert.NoError(t, err)
	assert.Equal(t, "value1", s1)

	s2, err := FetchKey(key2)
	assert.NoError(t, err)
	assert.Equal(t, "value2", s2)

	assert.NoError(t, DeleteKey(key1))
	assert.NotPanics(t, func() { MustDeleteKey(key2) })
}

func TestFetchKey(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(ctx, t)()

	s, err := FetchKey(<-random.ID + <-random.ID)
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestTTLKey(t *testing.T) {
	ctx := context.Background()
	defer mockery.Start(ctx, t)()

	k, v := <-random.ID, <-random.ID
	d, err := TTLKey(k)
	assert.Error(t, err)
	assert.Zero(t, d)

	assert.NoError(t, StoreKey(k, v, time.Minute))
	assert.NotPanics(t, func() { d = MustTTLKey(k) })
	assert.True(t, d <= time.Minute)
}
