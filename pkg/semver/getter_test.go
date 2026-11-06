package semver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/part"
	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

func TestSemverNew(t *testing.T) {
	v := semver.New(part.Uint64(1), part.Uint64(2), part.Uint64(3),
		part.NewParts("alpha"), "build")
	assert.Equal(t, "1.2.3-alpha+build", v.String())

	// wildcard
	v2 := semver.New(part.Any(true), part.Zero, part.Zero,
		nil, "")
	assert.True(t, v2.IsAny())
}

func TestSemverGetters(t *testing.T) {
	v, err := semver.Parse("1.2.3-alpha+build")
	require.NoError(t, err)

	assert.Equal(t, "1.2.3-alpha+build", v.String())
	assert.Equal(t, "1.2.3-alpha+build", v.Original())
	assert.Equal(t, "1.2.3+build", v.Release().String())
	assert.True(t, v.IsPreRelease())
	assert.Equal(t, "build", v.Metadata())

	// non-prerelease
	v2, _ := semver.Parse("1.2.3")
	assert.False(t, v2.IsPreRelease())
	assert.Equal(t, "", v2.Metadata())
}

func TestSemverMin(t *testing.T) {
	// wildcards must be constructed via New, not Parse
	v := semver.New(part.Uint64(1), part.Uint64(2), part.Any(true),
		nil, "")
	assert.Equal(t, "1.2.0", v.Min().String())

	v2 := semver.New(part.Uint64(1), part.Any(true), part.Any(true),
		nil, "")
	assert.Equal(t, "1.0.0", v2.Min().String())

	v3 := semver.New(part.Any(true), part.Any(true), part.Any(true),
		nil, "")
	assert.Equal(t, "0.0.0", v3.Min().String())

	// concrete version unchanged by Min
	v4, _ := semver.Parse("1.2.3")
	assert.Equal(t, "1.2.3", v4.Min().String())
}

func TestSemverRelease(t *testing.T) {
	v, _ := semver.Parse("1.2.3-alpha")
	assert.Equal(t, "1.2.3", v.Release().String())
}

func TestSemverCaretBump(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1.2.3", "2.0.0"},
		{"0.2.3", "0.3.0"},
		{"0.0.3", "0.0.4"},
		{"0.0.0", "0.0.1"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			v, err := semver.Parse(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, v.CaretBump().String())
		})
	}
}

func TestSemverTildeBump(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1.2.3", "1.3.0"},
		{"0.2.3", "0.3.0"},
		{"0.0.3", "0.1.0"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			v, err := semver.Parse(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, v.TildeBump().String())
		})
	}
}

func TestSemverInc(t *testing.T) {
	v, _ := semver.Parse("1.2.3")
	assert.Equal(t, "2.0.0", v.IncMajor().String())
	assert.Equal(t, "1.3.0", v.IncMinor().String())
	assert.Equal(t, "1.2.4", v.IncPatch().String())
}

func TestSemverComparison(t *testing.T) {
	a, _ := semver.Parse("1.0.0")
	b, _ := semver.Parse("2.0.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.True(t, a.LessThanOrEqual(b))
	assert.True(t, b.GreaterThanOrEqual(a))
	assert.False(t, a.Equal(b))

	aeq, _ := semver.Parse("1.0.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.Compare(aeq) == 0)
	assert.True(t, a.Compare(b) < 0)
	assert.True(t, b.Compare(a) > 0)
}
