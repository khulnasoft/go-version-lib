package gem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/gem"
)

func TestNewVersion(t *testing.T) {
	v, err := gem.NewVersion("1.2.3")
	require.NoError(t, err)
	assert.Equal(t, "1.2.3", v.String())

	// empty becomes "0"
	v, err = gem.NewVersion("")
	require.NoError(t, err)
	assert.Equal(t, "0", v.String())

	// invalid
	_, err = gem.NewVersion("not a version")
	require.Error(t, err)
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", -1},
		{"1.2.4", "1.2.3", 1},
		{"1.0", "2.0", -1},
		{"2.0", "1.0", 1},
		{"1.2.3", "1.2.3.4", -1},
		{"1.2.3.4", "1.2.3", 1},
		{"1.2.3-alpha", "1.2.3", -1},
		{"1.2.3", "1.2.3-alpha", 1},
	}

	for _, tc := range tests {
		t.Run(tc.v1+" vs "+tc.v2, func(t *testing.T) {
			a, err := gem.NewVersion(tc.v1)
			require.NoError(t, err)
			b, err := gem.NewVersion(tc.v2)
			require.NoError(t, err)
			assert.Equal(t, tc.want, a.Compare(b))
			assert.Equal(t, -tc.want, b.Compare(a))
		})
	}
}

func TestRelations(t *testing.T) {
	a, _ := gem.NewVersion("1.0")
	b, _ := gem.NewVersion("2.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.True(t, a.LessThanOrEqual(b))
	assert.True(t, b.GreaterThanOrEqual(a))
	assert.False(t, a.Equal(b))

	aeq, _ := gem.NewVersion("1.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.GreaterThanOrEqual(aeq))
	assert.True(t, a.LessThanOrEqual(aeq))
}

func TestRelease(t *testing.T) {
	v, _ := gem.NewVersion("1.2.0.a")
	rel := v.Release()
	assert.Equal(t, "1.2.0", rel.String())

	// non-prerelease returns itself
	v2, _ := gem.NewVersion("1.2.3")
	assert.Equal(t, "1.2.3", v2.Release().String())
}

func TestBump(t *testing.T) {
	v, _ := gem.NewVersion("5.3.1")
	assert.Equal(t, "5.4", v.Bump().String())

	v2, _ := gem.NewVersion("1.0")
	assert.Equal(t, "2", v2.Bump().String())
}
