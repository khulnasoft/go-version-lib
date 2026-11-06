package npm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/npm"
)

func TestConstraintNewConstraintsErrors(t *testing.T) {
	_, err := npm.NewConstraints(">= bar")
	require.Error(t, err)

	_, err = npm.NewConstraints("BAR >= 1.2.3")
	require.Error(t, err)
}

func TestConstraintString(t *testing.T) {
	c, err := npm.NewConstraints(">= 1.2.3, < 2.0 || > 3.0")
	require.NoError(t, err)
	assert.Contains(t, c.String(), ">= 1.2.3")
	assert.Contains(t, c.String(), "< 2.0")
	assert.Contains(t, c.String(), "||")
}

func TestConstraintCheck(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		{">= 1.0", "1.0.0", true},
		{">= 1.0", "0.9.0", false},
		{"< 2.0.0", "1.9.0", true},
		{"< 2.0.0", "2.0.0", false},
		{"1.0.0", "1.0.0", true},
		{"1.0.0", "1.0.1", false},
		{"*", "9.9.9", true},
		{"", "1.2.3", true},
		{"", "1.2.3-alpha.1", false},
	}

	for _, tc := range tests {
		t.Run(tc.constraint+"|"+tc.version, func(t *testing.T) {
			c, err := npm.NewConstraints(tc.constraint)
			require.NoError(t, err)
			v, err := npm.NewVersion(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Check(v))
		})
	}
}

func TestWithPreRelease(t *testing.T) {
	// Without WithPreRelease: a pre-release version is excluded from a
	// non-pre-release constraint.
	c, err := npm.NewConstraints(">= 1.0.0")
	require.NoError(t, err)
	v, err := npm.NewVersion("1.0.0-alpha")
	require.NoError(t, err)
	assert.False(t, c.Check(v))

	// With WithPreRelease: same pre-release satisfies itself.
	c2, err := npm.NewConstraints(">= 1.0.0-alpha", npm.WithPreRelease(true))
	require.NoError(t, err)
	v2, err := npm.NewVersion("1.0.0-alpha")
	require.NoError(t, err)
	assert.True(t, c2.Check(v2))
}

func TestWithBuildMetadata(t *testing.T) {
	// Without option: build metadata ignored -> equal
	c, err := npm.NewConstraints("=1.2.3")
	require.NoError(t, err)
	v, _ := npm.NewVersion("1.2.3+build.1")
	assert.True(t, c.Check(v))

	// With build metadata significant: "1.2.3+build.1" > "1.2.3" so != 1.2.3
	c2, err := npm.NewConstraints("=1.2.3", npm.WithBuildMetadata(true))
	require.NoError(t, err)
	assert.False(t, c2.Check(v))

	// but ">1.2.3" is satisfied by a version with build metadata
	c3, err := npm.NewConstraints(">1.2.3", npm.WithBuildMetadata(true))
	require.NoError(t, err)
	assert.True(t, c3.Check(v))
}

func TestVersionCollection(t *testing.T) {
	v1, _ := npm.NewVersion("1.0.0")
	v2, _ := npm.NewVersion("2.0.0")
	v3, _ := npm.NewVersion("3.0.0")
	col := npm.Collection{v2, v1, v3}
	assert.Equal(t, 3, col.Len())
	assert.True(t, col.Less(1, 0)) // 1 < 2
	col.Swap(0, 1)
	assert.Equal(t, "1.0.0", col[0].String())
}
