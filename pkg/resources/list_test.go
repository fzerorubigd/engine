package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	RegisterResource("a", "b")

	s, ok := QueryResource("a")
	assert.True(t, ok)
	assert.Equal(t, s, "b")

	s, ok = QueryResource("aa")
	assert.False(t, ok)
	assert.Empty(t, s)

}
