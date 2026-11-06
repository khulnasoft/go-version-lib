package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	debversion "github.com/khulnasoft/go-version-lib/pkg/deb"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		in   string
		err  bool
		want string
	}{
		{"1.0", false, "1.0"},
		{"1:1.0", false, "1:1.0"},
		{"1.0-2", false, "1.0-2"},
		{"1:1.0-2", false, "1:1.0-2"},
		{"1.0~rc1", false, "1.0~rc1"},
		{"", true, ""},
		{"abc1.0", true, ""}, // upstream must start with digit
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := debversion.NewVersion(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got.String())
		})
	}
}

func TestValid(t *testing.T) {
	assert.True(t, debversion.Valid("1.0"))
	assert.True(t, debversion.Valid("1:1.0-2"))
	assert.False(t, debversion.Valid(""))
	assert.False(t, debversion.Valid("abc1.0"))
}

func TestString(t *testing.T) {
	v, _ := debversion.NewVersion("1.0")
	assert.Equal(t, "1.0", v.String())

	v, _ = debversion.NewVersion("1:1.0")
	assert.Equal(t, "1:1.0", v.String())

	v, _ = debversion.NewVersion("1.0-2")
	assert.Equal(t, "1.0-2", v.String())

	v, _ = debversion.NewVersion("1:1.0-2")
	assert.Equal(t, "1:1.0-2", v.String())
}

func TestGetters(t *testing.T) {
	v, _ := debversion.NewVersion("1:1.0-2")
	assert.Equal(t, 1, v.Epoch())
	assert.Equal(t, "1.0", v.Version())
	assert.Equal(t, "2", v.Revision())
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		// equal
		{"1.0", "1.0", 0},
		{"1:1.0", "1:1.0", 0},
		{"1.0-2", "1.0-2", 0},
		// epoch dominates
		{"1:1.0", "0:9.9", 1},
		{"0:1.0", "1:1.0", -1},
		// upstream numeric compare
		{"1.0", "2.0", -1},
		{"2.0", "1.0", 1},
		{"1.10", "1.2", 1},
		// revision
		{"1.0-2", "1.0-1", 1},
		{"1.0-1", "1.0-2", -1},
		// tilde sorts before
		{"1.0~rc1", "1.0", -1},
		{"1.0", "1.0~rc1", 1},
		// plus sorts after letters
		{"1.0+dfsg", "1.0a", 1},
	}

	for _, tc := range tests {
		t.Run(tc.v1+" vs "+tc.v2, func(t *testing.T) {
			a, err := debversion.NewVersion(tc.v1)
			require.NoError(t, err)
			b, err := debversion.NewVersion(tc.v2)
			require.NoError(t, err)
			ab := a.Compare(b)
			ba := b.Compare(a)
			if tc.want == 0 {
				assert.Equal(t, 0, ab)
				assert.Equal(t, 0, ba)
			} else if tc.want < 0 {
				assert.Negative(t, ab)
				assert.Positive(t, ba)
			} else {
				assert.Positive(t, ab)
				assert.Negative(t, ba)
			}
		})
	}
}

func TestRelations(t *testing.T) {
	a, _ := debversion.NewVersion("1.0")
	b, _ := debversion.NewVersion("2.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.True(t, a.LessThanOrEqual(b))
	assert.True(t, b.GreaterThanOrEqual(a))
	assert.False(t, a.Equal(b))
	assert.False(t, a.GreaterThan(b))
	assert.False(t, b.LessThan(a))

	aeq, _ := debversion.NewVersion("1.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.GreaterThanOrEqual(aeq))
	assert.True(t, a.LessThanOrEqual(aeq))
}
