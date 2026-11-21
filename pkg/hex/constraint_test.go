package hex_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/hex"
)

func TestPessimistic(t *testing.T) {
	c, err := hex.NewConstraints("~> 1.2.3")
	require.NoError(t, err)
	in, _ := hex.NewVersion("1.2.9")
	out, _ := hex.NewVersion("1.3.0")
	assert.True(t, c.Check(in))
	assert.False(t, c.Check(out))
}

func TestPessimisticTwoSegment(t *testing.T) {
	c, err := hex.NewConstraints("~> 1.2")
	require.NoError(t, err)
	in, _ := hex.NewVersion("1.9.0")
	out, _ := hex.NewVersion("2.0.0")
	assert.True(t, c.Check(in))
	assert.False(t, c.Check(out))
}

func TestAndOr(t *testing.T) {
	c, err := hex.NewConstraints(">= 1.0.0 and < 2.0.0")
	require.NoError(t, err)
	v, _ := hex.NewVersion("1.5.0")
	assert.True(t, c.Check(v))
}
