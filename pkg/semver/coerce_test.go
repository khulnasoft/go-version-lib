package semver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

func TestCoerce(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  bool
	}{
		{"1.2.3", "1.2.3", false},
		{"v1.2", "1.2.0", false},
		{"1", "1.0.0", false},
		{"foo-1.2.3-bar", "1.2.3", false},
		{"release-2.0.0-rc1", "2.0.0", false},
		{"", "", true},
		{"no digits", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			v, err := semver.Coerce(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, v.String())
		})
	}
}
