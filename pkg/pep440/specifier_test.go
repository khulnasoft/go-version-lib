package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
)

func TestPep440NewSpecifiers(t *testing.T) {
	tests := []struct {
		in  string
		err bool
	}{
		{">= 1.0", false},
		{"~= 1.0", false},
		{"1.0", false},
		{"!= 1.0", false},
		{">= 1.0, < 2.0", false},
		{"1.0 || 2.0", false},
		{"improper [", true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			_, err := pep440version.NewSpecifiers(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPep440SpecifierCheck(t *testing.T) {
	tests := []struct {
		specifier string
		version   string
		want      bool
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
		t.Run(tc.specifier+"|"+tc.version, func(t *testing.T) {
			s, err := pep440version.NewSpecifiers(tc.specifier)
			require.NoError(t, err)
			v, err := pep440version.Parse(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, s.Check(v))
		})
	}
}

func TestPep440SpecifierCompatible(t *testing.T) {
	// ~= 1.2 means >= 1.2, < 2.0 (compatible release)
	s, err := pep440version.NewSpecifiers("~= 1.2")
	require.NoError(t, err)

	v12, _ := pep440version.Parse("1.2")
	assert.True(t, s.Check(v12))

	v19, _ := pep440version.Parse("1.9")
	assert.True(t, s.Check(v19))

	v20, _ := pep440version.Parse("2.0")
	assert.False(t, s.Check(v20))

	v11, _ := pep440version.Parse("1.1")
	assert.False(t, s.Check(v11))
}

func TestPep440SpecifierString(t *testing.T) {
	s, err := pep440version.NewSpecifiers(">= 1.0, < 2.0")
	require.NoError(t, err)
	assert.Contains(t, s.String(), ">= 1.0")
	assert.Contains(t, s.String(), "< 2.0")
}

func TestPep440SpecifierWithPreRelease(t *testing.T) {
	// Same pre-release satisfies itself regardless of the option.
	s, err := pep440version.NewSpecifiers(">= 1.0.0-alpha")
	require.NoError(t, err)
	v, _ := pep440version.Parse("1.0.0-alpha")
	assert.True(t, s.Check(v))

	// Without the option, a pre-release is excluded from a non-pre-release
	// constraint.
	s2, err := pep440version.NewSpecifiers(">= 1.0.0")
	require.NoError(t, err)
	assert.False(t, s2.Check(v))

	// With the option, a pre-release can satisfy a non-pre-release constraint
	// when precedence allows (e.g. equal pre-release).
	s3, err := pep440version.NewSpecifiers(">= 1.0.0-alpha", pep440version.WithPreRelease(true))
	require.NoError(t, err)
	assert.True(t, s3.Check(v))
}

func TestPep440SpecifierCompare(t *testing.T) {
	// Version comparison methods
	a, _ := pep440version.Parse("1.0")
	b, _ := pep440version.Parse("2.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.False(t, a.Equal(b))

	aeq, _ := pep440version.Parse("1.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.Compare(aeq) == 0)
}
