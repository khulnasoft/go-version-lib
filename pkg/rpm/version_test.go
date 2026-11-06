package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	rpmversion "github.com/khulnasoft/go-version-lib/pkg/rpm"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1.0", "1.0"},
		{"1:1.0", "1:1.0"},
		{"1.0-2", "1.0-2"},
		{"1:1.0-2", "1:1.0-2"},
		{"1.0~rc1", "1.0~rc1"},
		{"1.0^git1", "1.0^git1"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			v := rpmversion.NewVersion(tc.in)
			assert.Equal(t, tc.want, v.String())
		})
	}
}

func TestGetters(t *testing.T) {
	v := rpmversion.NewVersion("1:1.0-2")
	assert.Equal(t, 1, v.Epoch())
	assert.Equal(t, "1.0", v.Version())
	assert.Equal(t, "2", v.Release())
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		// equal
		{"1.0", "1.0", 0},
		{"1:1.0-2", "1:1.0-2", 0},
		// epoch dominates
		{"1:1.0", "0:9.9", 1},
		{"0:1.0", "1:1.0", -1},
		// numeric compare, longer wins
		{"1.10", "1.2", 1},
		{"1.2", "1.10", -1},
		// leading zeros ignored
		{"1.01", "1.1", 0},
		// number > alpha (at the same segment position)
		{"1.0", "1a", 1},
		{"1a", "1.0", -1},
		// tilde sorts before everything
		{"1.0~rc1", "1.0", -1},
		{"1.0", "1.0~rc1", 1},
		{"1.0~rc1", "1.0~rc2", -1},
		// caret sorts after base version
		{"1.0", "1.0^git1", -1},
		{"1.0^git1", "1.0", 1},
		{"1.0^git1", "1.0.1", -1},
		// release compare
		{"1.0-2", "1.0-1", 1},
		{"1.0-1", "1.0-2", -1},
		// most segments wins
		{"1.0.1", "1.0", 1},
	}

	for _, tc := range tests {
		t.Run(tc.v1+" vs "+tc.v2, func(t *testing.T) {
			a := rpmversion.NewVersion(tc.v1)
			b := rpmversion.NewVersion(tc.v2)
			assert.Equal(t, tc.want, a.Compare(b))
			assert.Equal(t, -tc.want, b.Compare(a))
		})
	}
}

func TestRelations(t *testing.T) {
	a := rpmversion.NewVersion("1.0")
	b := rpmversion.NewVersion("2.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.True(t, a.LessThanOrEqual(b))
	assert.True(t, b.GreaterThanOrEqual(a))
	assert.False(t, a.Equal(b))
	assert.False(t, a.GreaterThan(b))
	assert.False(t, b.LessThan(a))

	aeq := rpmversion.NewVersion("1.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.GreaterThanOrEqual(aeq))
	assert.True(t, a.LessThanOrEqual(aeq))
}
