package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/version"
)

func TestBump(t *testing.T) {
	cases := []struct {
		in   string
		kind version.BumpKind
		want string
	}{
		{"1.2.3", version.BumpMajor, "2.0.0"},
		{"1.2.3", version.BumpMinor, "1.3.0"},
		{"1.2.3", version.BumpPatch, "1.2.4"},
		{"1.2", version.BumpPatch, "1.2.1"},
		{"1.2.3-alpha+build", version.BumpMinor, "1.3.0"},
		{"1.2.3.4", version.BumpMajor, "2.0.0.0"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			v, err := version.Parse(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, v.Bump(tc.kind).String())
		})
	}
}
