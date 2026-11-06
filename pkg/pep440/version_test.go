package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
)

func TestPep440Getters(t *testing.T) {
	v, err := pep440version.Parse("2!1.0a1.post1.dev1+local.1")
	require.NoError(t, err)

	assert.Equal(t, "2!1.0a1.post1.dev1+local.1", v.String())
	assert.Equal(t, "2!1.0", v.BaseVersion())
	assert.Equal(t, "local.1", v.Local())
	assert.Equal(t, "2!1.0a1.post1.dev1", v.Public())
	assert.True(t, v.IsPreRelease())
	assert.True(t, v.IsPostRelease())

	// release version
	v2, _ := pep440version.Parse("1.0")
	assert.False(t, v2.IsPreRelease())
	assert.False(t, v2.IsPostRelease())
	assert.Equal(t, "1.0", v2.BaseVersion())
	assert.Equal(t, "", v2.Local())
	assert.Equal(t, "1.0", v2.Public())
}

func TestPep440MarshalUnmarshalText(t *testing.T) {
	v, err := pep440version.Parse("1.0a1.post1+local")
	require.NoError(t, err)

	data, err := v.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "1.0a1.post1+local", string(data))

	var v2 pep440version.Version
	require.NoError(t, v2.UnmarshalText(data))
	assert.Equal(t, v.String(), v2.String())
}

func TestPep440MustParse(t *testing.T) {
	assert.NotPanics(t, func() { pep440version.MustParse("1.0") })
	assert.Panics(t, func() { pep440version.MustParse("not-a-version") })
}

func TestPep440Original(t *testing.T) {
	v, err := pep440version.Parse("  v1.0  ")
	require.NoError(t, err)
	assert.Equal(t, "  v1.0  ", v.Original())
}
