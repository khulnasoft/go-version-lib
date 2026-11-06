package reference_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/reference"
)

func TestDetectPurl(t *testing.T) {
	tests := []struct {
		in   string
		want reference.VersionType
		err  bool
	}{
		{"pkg:npm/lodash@4.17.21", reference.NPM, false},
		{"pkg:deb/foo/bar@1.0-2", reference.Deb, false},
		{"pkg:rpm/foo@1.0", reference.RPM, false},
		{"pkg:maven/a/b@1.0", reference.Maven, false},
		{"pkg:gem/foo@1.2.3", reference.Gem, false},
		{"pkg:pep440/foo@1.0a1", reference.PEP440, false},
		{"pkg:semver/foo@1.2.3", reference.Semver, false},
		{"pkg:unknown/foo@1.0", reference.Unknown, true},
		{"pkg:", reference.Unknown, true},
		{"foo bar", reference.Unknown, true},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := reference.Detect(tc.in)
			if tc.err {
				require.Error(t, err)
				assert.Equal(t, reference.Unknown, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewReferencePurl(t *testing.T) {
	ref, err := reference.NewReference("pkg:npm/lodash@4.17.21")
	require.NoError(t, err)
	assert.Equal(t, reference.NPM, ref.Type())
	assert.Equal(t, "pkg:npm/lodash@4.17.21", ref.String())

	// Round trip: the parsed underlying version is comparable.
	same, err := reference.NewReference("pkg:npm/lodash@4.17.21")
	require.NoError(t, err)
	assert.True(t, ref.Equal(same))
	assert.Equal(t, 0, ref.Compare(same))

	diff, err := reference.NewReference("pkg:npm/lodash@1.0.0")
	require.NoError(t, err)
	assert.False(t, ref.Equal(diff))
	assert.True(t, ref.GreaterThan(diff))
	assert.True(t, diff.LessThan(ref))

	// Different package types are not equal.
	rpm, err := reference.NewReference("pkg:rpm/lodash@4.17.21")
	require.NoError(t, err)
	assert.False(t, ref.Equal(rpm))

	// A purl without a version is an error.
	_, err = reference.NewReference("pkg:npm/lodash")
	require.Error(t, err)

	// An unknown package type is an error.
	_, err = reference.NewReference("pkg:unknown/foo@1.0")
	require.Error(t, err)
}

func TestNewReferencePurlPreservesPackageIdentity(t *testing.T) {
	a, err := reference.NewReference("pkg:npm/foo@1.0.0")
	require.NoError(t, err)
	b, err := reference.NewReference("pkg:npm/bar@1.0.0")
	require.NoError(t, err)
	// Same version, different packages: not equal as package references.
	assert.False(t, a.Equal(b))
	assert.NotEqual(t, a.String(), b.String())
}

func TestNewReferencePurlVersionOrdering(t *testing.T) {
	low, err := reference.NewReference("pkg:npm/foo@1.2.0")
	require.NoError(t, err)
	high, err := reference.NewReference("pkg:npm/foo@1.10.0")
	require.NoError(t, err)
	// Same package: versions compared version-aware, so 1.10.0 > 1.2.0
	// (lexicographic string comparison would get this wrong).
	assert.True(t, high.GreaterThan(low))
	assert.True(t, low.LessThan(high))
	assert.False(t, high.Equal(low))
}

func TestNewReferencePlainVersion(t *testing.T) {
	// Plain version strings are detected and parsed into ecosystem refs.
	// Detection is order-based (semver first), so overlapping regexes resolve
	// to the first matching type.
	semverRef, err := reference.NewReference("1.2.3")
	require.NoError(t, err)
	assert.Equal(t, reference.Semver, semverRef.Type())
	assert.Equal(t, "1.2.3", semverRef.String())

	// Deb is detectable via epoch syntax.
	debRef, err := reference.NewReference("1:1.0-2")
	require.NoError(t, err)
	assert.Equal(t, reference.Deb, debRef.Type())
	assert.Equal(t, "1:1.0-2", debRef.String())

	// Gem is detectable via dotted-alpha segment.
	gemRef, err := reference.NewReference("1.0.a")
	require.NoError(t, err)
	assert.Equal(t, reference.Gem, gemRef.Type())
	assert.Equal(t, "1.0.a", gemRef.String())

	// PEP440 is detectable via epoch bang (!) at the start, which gem rejects.
	pep440Ref, err := reference.NewReference("2!3.4b1")
	require.NoError(t, err)
	assert.Equal(t, reference.PEP440, pep440Ref.Type())
	assert.Equal(t, "2!3.4b1", pep440Ref.String())

	// Maven is detectable via NewReferenceByType.
	mvnRef, err := reference.NewReferenceByType("1.0", reference.Maven)
	require.NoError(t, err)
	assert.Equal(t, reference.Maven, mvnRef.Type())
	assert.Equal(t, "1.0", mvnRef.String())

	// Ref methods all work within the same type.
	assert.True(t, semverRef.GreaterThanOrEqual(semverRef))
	assert.True(t, semverRef.LessThanOrEqual(semverRef))
	assert.False(t, semverRef.GreaterThan(semverRef))
	assert.False(t, semverRef.LessThan(semverRef))
	assert.True(t, semverRef.Equal(semverRef))
	assert.Equal(t, 0, semverRef.Compare(semverRef))

	// GreaterThan/LessThan work for same-type refs.
	low, _ := reference.NewReference("1.0.0")
	high, _ := reference.NewReference("2.0.0")
	assert.True(t, low.LessThan(high))
	assert.True(t, high.GreaterThan(low))
	assert.False(t, low.Equal(high))
}

func TestNewReferenceByType(t *testing.T) {
	v, err := reference.NewReferenceByType("1.2.3", reference.Semver)
	require.NoError(t, err)
	assert.Equal(t, reference.Semver, v.Type())
	assert.Equal(t, "1.2.3", v.String())

	_, err = reference.NewReferenceByType("1.2.3", reference.Unknown)
	require.Error(t, err)
}

func TestDetectPlainVersion(t *testing.T) {
	tests := []struct {
		in   string
		want reference.VersionType
	}{
		{"1.2.3", reference.Semver},
		{"1:1.0", reference.Deb},
		{"1.0.a", reference.Gem},
		{"2!3.4b1", reference.PEP440},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := reference.Detect(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRefTypeAllMethods(t *testing.T) {
	// Exercise every Reference interface method on each ecosystem ref type.
	cases := []struct {
		in   string
		want reference.VersionType
	}{
		{"1.2.3", reference.Semver},
		{"1.0.a", reference.Gem},
		{"2!3.4b1", reference.PEP440},
		{"1:1.0-2", reference.Deb},
		{"1.0^1", reference.RPM},
	}

	for _, tc := range cases {
		t.Run(string(tc.want), func(t *testing.T) {
			a, err := reference.NewReference(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, a.Type())

			b, err := reference.NewReference(tc.in)
			require.NoError(t, err)

			// All comparison methods
			assert.True(t, a.Equal(b))
			assert.Equal(t, 0, a.Compare(b))
			assert.True(t, a.GreaterThanOrEqual(b))
			assert.True(t, a.LessThanOrEqual(b))
			assert.False(t, a.GreaterThan(b))
			assert.False(t, a.LessThan(b))

			// String round-trips
			assert.Equal(t, tc.in, a.String())

			// A different version compares correctly
			c, err := reference.NewReference("9.9.9")
			require.NoError(t, err)
			if c.Type() == tc.want {
				assert.True(t, c.GreaterThan(a))
				assert.True(t, a.LessThan(c))
				assert.False(t, c.Equal(a))
			}
		})
	}

	// NPM and Maven are tested via NewReferenceByType since plain-string
	// detection resolves them to earlier-registered types.
	npmRef, err := reference.NewReferenceByType("1.2.3", reference.NPM)
	require.NoError(t, err)
	assert.Equal(t, reference.NPM, npmRef.Type())
	assert.True(t, npmRef.GreaterThanOrEqual(npmRef))
	assert.True(t, npmRef.LessThanOrEqual(npmRef))
	assert.Equal(t, "1.2.3", npmRef.String())

	mvnRef, err := reference.NewReferenceByType("1.0", reference.Maven)
	require.NoError(t, err)
	assert.Equal(t, reference.Maven, mvnRef.Type())
	assert.True(t, mvnRef.GreaterThanOrEqual(mvnRef))
	assert.True(t, mvnRef.LessThanOrEqual(mvnRef))
	assert.Equal(t, "1.0", mvnRef.String())
}
