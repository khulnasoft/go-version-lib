package npm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/npm"
)

func TestToVersionRangeChecksParity(t *testing.T) {
	constraints := []string{
		">=1.2.3",
		"<2.0.0",
		">=1.0.0 <2.0.0",
		"1.2.3",
		"^1.2.3",
		"~1.2.3",
		">=1.0.0 || <0.5.0",
		"1.x",
		"1.2.x",
		"*",
	}
	versions := []string{"0.4.0", "0.5.0", "1.0.0", "1.2.3", "1.2.4", "1.3.0", "1.9.9", "2.0.0", "2.5.0"}

	for _, cs := range constraints {
		t.Run(cs, func(t *testing.T) {
			c, err := npm.NewConstraints(cs)
			require.NoError(t, err)
			r, err := c.ToVersionRange()
			require.NoError(t, err)

			for _, vs := range versions {
				v, err := npm.NewVersion(vs)
				require.NoError(t, err)
				assert.Equal(t, c.Check(v), r.Satisfies(v),
					"constraint %q vs version %s", cs, vs)
			}
		})
	}
}
