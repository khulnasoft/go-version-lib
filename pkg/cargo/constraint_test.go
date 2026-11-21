package cargo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/cargo"
)

func TestBareVersionIsCaret(t *testing.T) {
	c, err := cargo.NewConstraints("1.2.3")
	require.NoError(t, err)
	// caret: >=1.2.3 <2.0.0
	v129, _ := cargo.NewVersion("1.2.9")
	v200, _ := cargo.NewVersion("2.0.0")
	v122, _ := cargo.NewVersion("1.2.2")
	assert.True(t, c.Check(v129))
	assert.False(t, c.Check(v200))
	assert.False(t, c.Check(v122))
}

func TestWildcard(t *testing.T) {
	c, err := cargo.NewConstraints("1.2.*")
	require.NoError(t, err)
	in, _ := cargo.NewVersion("1.2.9")
	out, _ := cargo.NewVersion("1.3.0")
	assert.True(t, c.Check(in))
	assert.False(t, c.Check(out))
}

func TestExplicitExact(t *testing.T) {
	c, err := cargo.NewConstraints("=1.2.3")
	require.NoError(t, err)
	v, _ := cargo.NewVersion("1.2.4")
	assert.False(t, c.Check(v))
}
