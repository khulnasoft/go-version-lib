package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mvnversion "github.com/khulnasoft/go-version-lib/pkg/mvn"
)

func TestToVersionRangeChecksParity(t *testing.T) {
	constraints := []string{
		">=1.2.3",
		"<2.0.0",
		">=1.0.0 <2.0.0",
		"1.2.3",
		"!=1.5.0",
		">=1.0.0 || <0.5.0",
	}
	versions := []string{"0.4.0", "0.5.0", "1.0.0", "1.2.3", "1.2.4", "1.3.0", "1.9.9", "1.5.0", "2.0.0", "2.5.0"}

	for _, cs := range constraints {
		t.Run(cs, func(t *testing.T) {
			c, err := mvnversion.NewConstraints(cs)
			require.NoError(t, err)
			r, err := c.ToVersionRange()
			require.NoError(t, err)

			for _, vs := range versions {
				v, err := mvnversion.NewVersion(vs)
				require.NoError(t, err)
				assert.Equal(t, c.Check(v), r.Satisfies(v),
					"constraint %q vs version %s", cs, vs)
			}
		})
	}
}

// Maven's ordering is non-transitive for qualifier versions; bound merging
// must not drop predicates that Check enforces.
func TestToVersionRangeNonTransitiveBounds(t *testing.T) {
	c, err := mvnversion.NewConstraints(">1-rc >1.sp")
	require.NoError(t, err)
	r, err := c.ToVersionRange()
	require.NoError(t, err)

	// 1.sp.1 is >1.sp but <1-rc: Check rejects it, so the range must too
	for _, vs := range []string{"1.sp.1", "1-rc.1", "2.0"} {
		v, err := mvnversion.NewVersion(vs)
		require.NoError(t, err)
		assert.Equal(t, c.Check(v), r.Satisfies(v), vs)
	}
}
