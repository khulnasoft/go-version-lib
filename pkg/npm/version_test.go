package npm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/npm"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
		err  bool
	}{
		{"1.2.3", "1.2.3", false},
		{"1.2.3-alpha", "1.2.3-alpha", false},
		{"1.2.3+build", "1.2.3+build", false},
		{"1.2.3-alpha.1+build.2", "1.2.3-alpha.1+build.2", false},
		{"v1.2.3", "", true}, // semver does not accept a leading "v"
		{"1", "", true},
		{"1.0", "", true},
		{"foo", "", true},
		{"1.2.3.4", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := npm.NewVersion(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got.String())
		})
	}
}

func TestCompare(t *testing.T) {
	a, err := npm.NewVersion("1.2.3")
	require.NoError(t, err)
	b, err := npm.NewVersion("1.2.4")
	require.NoError(t, err)

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.False(t, a.Equal(b))

	// build metadata ignored in equality
	c, _ := npm.NewVersion("1.2.3+build1")
	d, _ := npm.NewVersion("1.2.3+build2")
	assert.True(t, c.Equal(d))
}
