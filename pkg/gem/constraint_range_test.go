package gem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/gem"
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
			c, err := gem.NewConstraints(cs)
			require.NoError(t, err)
			r, err := c.ToVersionRange()
			require.NoError(t, err)

			for _, vs := range versions {
				v, err := gem.NewVersion(vs)
				require.NoError(t, err)
				assert.Equal(t, c.Check(v), r.Satisfies(v),
					"constraint %q vs version %s", cs, vs)
			}
		})
	}
}

func TestToVersionRangePessimisticError(t *testing.T) {
	c, err := gem.NewConstraints("~>1.2.3")
	require.NoError(t, err)
	_, err = c.ToVersionRange()
	assert.Error(t, err)
}
