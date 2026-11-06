package purl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/purl"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want purl.PackageURL
		err  bool
	}{
		{
			name: "simple",
			in:   "pkg:npm/lodash@4.17.21",
			want: purl.PackageURL{
				Type:    "npm",
				Name:    "lodash",
				Version: "4.17.21",
			},
		},
		{
			name: "with namespace",
			in:   "pkg:deb/foo/bar@1.0-2",
			want: purl.PackageURL{
				Type:      "deb",
				Namespace: "foo",
				Name:      "bar",
				Version:   "1.0-2",
			},
		},
		{
			name: "multi-level namespace",
			in:   "pkg:maven/a/b/c@1.0",
			want: purl.PackageURL{
				Type:      "maven",
				Namespace: "a/b",
				Name:      "c",
				Version:   "1.0",
			},
		},
		{
			name: "scoped name round trip",
			in:   "pkg:npm/@angular/core@1.0.0",
			want: purl.PackageURL{
				Type:      "npm",
				Namespace: "@angular",
				Name:      "core",
				Version:   "1.0.0",
			},
		},
		{
			name: "no version",
			in:   "pkg:npm/lodash",
			want: purl.PackageURL{
				Type: "npm",
				Name: "lodash",
			},
		},
		{
			name: "with qualifiers",
			in:   "pkg:rpm/foo@1.0?arch=x86_64&distro=fc41",
			want: purl.PackageURL{
				Type:    "rpm",
				Name:    "foo",
				Version: "1.0",
				Qualifiers: map[string]string{
					"arch":   "x86_64",
					"distro": "fc41",
				},
			},
		},
		{
			name: "with subpath",
			in:   "pkg:deb/foo@1.0#usr/bin/foo",
			want: purl.PackageURL{
				Type:    "deb",
				Name:    "foo",
				Version: "1.0",
				Subpath: "usr/bin/foo",
			},
		},
		{
			name: "full purl",
			in:   "pkg:deb/namespace/name@1.0-2?arch=amd64#sub/path",
			want: purl.PackageURL{
				Type:      "deb",
				Namespace: "namespace",
				Name:      "name",
				Version:   "1.0-2",
				Qualifiers: map[string]string{
					"arch": "amd64",
				},
				Subpath: "sub/path",
			},
		},
		{
			name: "pep440 epoch version",
			in:   "pkg:pep440/foo@2!3.4b1+local",
			want: purl.PackageURL{
				Type:    "pep440",
				Name:    "foo",
				Version: "2!3.4b1+local",
			},
		},
		{
			name: "url encoded namespace",
			in:   "pkg:deb/my%2Fspace/name@1.0",
			want: purl.PackageURL{
				Type:      "deb",
				Namespace: "my/space",
				Name:      "name",
				Version:   "1.0",
			},
		},
		{
			name: "url encoded name",
			in:   "pkg:npm/name%40scope@1.0",
			want: purl.PackageURL{
				Type:    "npm",
				Name:    "name@scope",
				Version: "1.0",
			},
		},
		{
			name: "missing prefix",
			in:   "npm/lodash@1.0",
			err:  true,
		},
		{
			name: "missing type",
			in:   "pkg://name@1.0",
			err:  true,
		},
		{
			name: "missing name",
			in:   "pkg:npm@1.0",
			err:  true,
		},
		{
			name: "empty",
			in:   "",
			err:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := purl.Parse(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStringRoundTrip(t *testing.T) {
	inputs := []string{
		"pkg:npm/lodash@4.17.21",
		"pkg:deb/foo/bar@1.0-2",
		"pkg:maven/a/b/c@1.0?foo=bar",
		"pkg:npm/@angular/core@1.0.0#sub",
		"pkg:rpm/name@1.0?arch=x86_64&distro=fc41#usr/bin",
		"pkg:gem/foo@1.2.3",
		"pkg:semver/foo@1.2.3",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			p, err := purl.Parse(in)
			require.NoError(t, err)
			assert.Equal(t, in, p.String())
		})
	}
}

func TestStringEmptyQualifiers(t *testing.T) {
	p, err := purl.Parse("pkg:npm/lodash@1.0")
	require.NoError(t, err)
	assert.Equal(t, "pkg:npm/lodash@1.0", p.String())
}

// TestEncodedRoundTrip verifies that percent-encoded characters decode on
// Parse and re-encode on String, forming a stable Parse/String cycle even
// when the canonical form is not byte-identical to the raw input.
func TestEncodedRoundTrip(t *testing.T) {
	// '!' and '+' are left unencoded in the canonical form.
	p, err := purl.Parse("pkg:pep440/foo@2!3.4b1+local")
	require.NoError(t, err)
	assert.Equal(t, "2!3.4b1+local", p.Version)
	assert.Equal(t, "pkg:pep440/foo@2!3.4b1+local", p.String())

	// Feeding the canonical form back in yields the same version.
	p2, err := purl.Parse(p.String())
	require.NoError(t, err)
	assert.Equal(t, p.Version, p2.Version)
	assert.Equal(t, p.String(), p2.String())

	// A literal '%' is encoded so it round-trips safely.
	p3, err := purl.Parse("pkg:deb/name@100%25pure")
	require.NoError(t, err)
	assert.Equal(t, "100%pure", p3.Version)
	assert.Equal(t, "pkg:deb/name@100%25pure", p3.String())

	// Encoded namespace/name separators decode back.
	p4, err := purl.Parse("pkg:deb/my%2Fspace/name%40scope@1.0")
	require.NoError(t, err)
	assert.Equal(t, "my/space", p4.Namespace)
	assert.Equal(t, "name@scope", p4.Name)
	assert.Equal(t, "1.0", p4.Version)
}
