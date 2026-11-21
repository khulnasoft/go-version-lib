package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/version"
)

// TestToVersionRangeChecksParity verifies that the lowered range agrees with
// Constraints.Check for a matrix of constraints and versions.
func TestToVersionRangeChecksParity(t *testing.T) {
	constraints := []string{
		">=1.2.3",
		"<2.0.0",
		">=1.0.0 <2.0.0",
		"1.2.3",
		"!=1.5.0",
		"^1.2.3",
		"~1.2.3",
		"~>1.2.3",
		">=1.0.0 || <0.5.0",
	}
	versions := []string{"0.4.0", "0.5.0", "1.0.0", "1.2.3", "1.2.4", "1.3.0", "1.9.9", "1.5.0", "2.0.0", "2.5.0"}

	for _, cs := range constraints {
		t.Run(cs, func(t *testing.T) {
			c, err := version.NewConstraints(cs)
			require.NoError(t, err)
			r, err := c.ToVersionRange()
			require.NoError(t, err)

			for _, vs := range versions {
				v, err := version.Parse(vs)
				require.NoError(t, err)
				assert.Equal(t, c.Check(v), r.Satisfies(v),
					"constraint %q vs version %s", cs, vs)
			}
		})
	}
}

func TestToVersionRangeIntersect(t *testing.T) {
	a, err := version.NewConstraints(">=1.0.0 <2.0.0")
	require.NoError(t, err)
	b, err := version.NewConstraints("^1.2.0")
	require.NoError(t, err)

	ra, err := a.ToVersionRange()
	require.NoError(t, err)
	rb, err := b.ToVersionRange()
	require.NoError(t, err)

	inter := ra.Intersect(rb)
	// intersection of [1.0,2.0) and [1.2.0,2.0.0) is [1.2.0,2.0.0)
	for _, tc := range []struct {
		v    string
		want bool
	}{
		{"1.1.9", false},
		{"1.2.0", true},
		{"1.5.0", true},
		{"2.0.0", false},
	} {
		v, err := version.Parse(tc.v)
		require.NoError(t, err)
		assert.Equal(t, tc.want, inter.Satisfies(v), tc.v)
	}
}
