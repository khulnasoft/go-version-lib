package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
)

// TestPep440Ordering exercises canonical PEP 440 ordering rules:
// epoch > release > pre-release > post/dev releases, with dev releases
// sorting before their associated pre/normal release.
func TestPep440Ordering(t *testing.T) {
	cases := []struct {
		lower string
		upper string
	}{
		// epochs dominate everything else
		{"2.0", "1!0.1"},
		{"0.9", "1!0.1"},
		// pre-releases sort before the final release
		{"1.0a1", "1.0"},
		{"1.0b2", "1.0rc1"},
		{"1.0rc1", "1.0"},
		// dev releases sort before pre-releases and finals
		{"1.0.dev1", "1.0a1"},
		{"1.0a1.dev1", "1.0a1"},
		// post releases sort after the final release
		{"1.0", "1.0.post1"},
		{"1.0.post1", "1.0.post2"},
		// local versions don't affect ordering vs non-local
		{"1.0", "1.0+local.1"},
	}
	for _, tc := range cases {
		t.Run(tc.lower+"<"+tc.upper, func(t *testing.T) {
			lo, err := pep440version.Parse(tc.lower)
			require.NoError(t, err)
			hi, err := pep440version.Parse(tc.upper)
			require.NoError(t, err)
			assert.Negative(t, lo.Compare(hi), "%s should be < %s", tc.lower, tc.upper)
			assert.Positive(t, hi.Compare(lo), "%s should be > %s", tc.upper, tc.lower)
		})
	}
}

// TestPep440Normalization checks that equivalent spellings compare equal.
func TestPep440Normalization(t *testing.T) {
	pairs := [][2]string{
		{"1.0", "1.0.0"},
		{"1.0a1", "1.0.a1"},
		{"1.0-a1", "1.0a1"},
		{"v1.2.3", "1.2.3"},
	}
	for _, p := range pairs {
		t.Run(p[0]+"=="+p[1], func(t *testing.T) {
			a, err := pep440version.Parse(p[0])
			require.NoError(t, err)
			b, err := pep440version.Parse(p[1])
			require.NoError(t, err)
			assert.Equal(t, 0, a.Compare(b), "%s should equal %s", p[0], p[1])
			assert.Equal(t, 0, b.Compare(a), "%s should equal %s (reverse)", p[1], p[0])
		})
	}
}
