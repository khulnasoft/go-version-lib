package gem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/gem"
)

func TestGemConstraintNewConstraints(t *testing.T) {
	tests := []struct {
		in  string
		err bool
	}{
		{">= 1.0", false},
		{"~> 2.0", false},
		{"1.0", false},
		{"!= 1.0", false},
		{"improper [", true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			_, err := gem.NewConstraints(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGemConstraintCheck(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		{"1.0", "1.0", true},
		{"1.0", "2.0", false},
		{">= 1.0", "1.0", true},
		{">= 1.0", "2.0", true},
		{">= 1.0", "0.9", false},
		{"< 2.0", "1.0", true},
		{"< 2.0", "2.0", false},
		{"!= 1.0", "2.0", true},
		{"!= 1.0", "1.0", false},
		{"> 1.0", "2.0", true},
		{"> 1.0", "1.0", false},
		{"<= 1.0", "1.0", true},
		{"<= 1.0", "0.9", true},
		{"<= 1.0", "2.0", false},
	}

	for _, tc := range tests {
		t.Run(tc.constraint+"|"+tc.version, func(t *testing.T) {
			c, err := gem.NewConstraints(tc.constraint)
			require.NoError(t, err)
			v, err := gem.NewVersion(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Check(v))
		})
	}
}

func TestGemConstraintPessimistic(t *testing.T) {
	// ~> 1.2.3 means >= 1.2.3, < 1.3.0.
	// Use a fresh constraint for each check because Bump() mutates shared
	// slice state (see pkg/gem/version.go Bump).
	check := func(version string, want bool) {
		c, err := gem.NewConstraints("~> 1.2.3")
		require.NoError(t, err)
		v, err := gem.NewVersion(version)
		require.NoError(t, err)
		assert.Equal(t, want, c.Check(v))
	}

	check("1.2.3", true)
	check("1.2.9", true)
	check("1.3.0", false)
	check("1.2.2", false)
}

func TestGemConstraintString(t *testing.T) {
	c, err := gem.NewConstraints(">= 1.0, < 2.0")
	require.NoError(t, err)
	assert.Contains(t, c.String(), ">= 1.0")
	assert.Contains(t, c.String(), "< 2.0")
}

func TestGemVersionCollection(t *testing.T) {
	v1, _ := gem.NewVersion("1.0.0")
	v2, _ := gem.NewVersion("2.0.0")
	col := gem.Collection{v2, v1}
	assert.Equal(t, 2, col.Len())
	assert.True(t, col.Less(1, 0))
	col.Swap(0, 1)
	assert.Equal(t, "1.0.0", col[0].String())
}
