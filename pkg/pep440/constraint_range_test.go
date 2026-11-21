package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
)

func TestToVersionRangeChecksParity(t *testing.T) {
	constraints := []string{
		">=1.2.3",
		"<2.0.0",
		">=1.0.0,<2.0.0",
		"==1.2.3",
		"!=1.5.0",
		"~=1.4.5",
		"~=1.4",
		"==1.4.*",
		">=1.0,<2.0 || <0.5",
	}
	versions := []string{"0.4.0", "0.5.0", "1.0", "1.2.3", "1.2.4", "1.3.0", "1.4.4", "1.4.5", "1.4.9", "1.5.0", "1.9.9", "1.4a1", "1.4.dev0", "1.4rc1", "1.4.post1", "2.0.0a1", "2.0.0.dev0", "2.0.0", "2.5.0"}

	for _, cs := range constraints {
		t.Run(cs, func(t *testing.T) {
			c, err := pep440version.NewSpecifiers(cs)
			require.NoError(t, err)
			r, err := c.ToVersionRange()
			require.NoError(t, err)

			for _, vs := range versions {
				v, err := pep440version.Parse(vs)
				require.NoError(t, err)
				assert.Equal(t, c.Check(v), r.Satisfies(v),
					"specifier %q vs version %s", cs, vs)
			}
		})
	}
}

func TestToVersionRangeErrors(t *testing.T) {
	for _, cs := range []string{"===1.2.3", "~=1.0a1", "!=1.4.*"} {
		c, err := pep440version.NewSpecifiers(cs)
		require.NoError(t, err, cs)
		_, err = c.ToVersionRange()
		assert.Error(t, err, cs)
	}
}
