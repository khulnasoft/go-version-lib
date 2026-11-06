package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mvnversion "github.com/khulnasoft/go-version-lib/pkg/mvn"
)

func TestNewVersion(t *testing.T) {
	v, err := mvnversion.NewVersion("1.0")
	require.NoError(t, err)
	assert.Equal(t, "1.0", v.String())
	assert.Equal(t, "1.0", v.Value)
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0", "1.0", 0},
		{"2.0", "1.0", 1},
		{"1.0", "2.0", -1},
		{"1.10", "1.2", 1},
		{"1.2", "1.10", -1},
		// release (empty) is highest among qualifiers
		{"1.0", "1.0-rc", 1},
		{"1.0-rc", "1.0", -1},
		// aliases
		{"1.0-ga", "1.0", 0},
		{"1.0-release", "1.0", 0},
		{"1.0-final", "1.0", 0},
		// qualifier ordering
		{"1.0-alpha", "1.0-beta", -1},
		{"1.0-beta", "1.0-milestone", -1},
		{"1.0-milestone", "1.0-rc", -1},
		{"1.0-rc", "1.0", -1},
		// snapshot < release
		{"1.0-SNAPSHOT", "1.0", -1},
		// cr alias -> rc
		{"1.0-cr", "1.0-rc", 0},
	}

	for _, tc := range tests {
		t.Run(tc.v1+" vs "+tc.v2, func(t *testing.T) {
			a, err := mvnversion.NewVersion(tc.v1)
			require.NoError(t, err)
			b, err := mvnversion.NewVersion(tc.v2)
			require.NoError(t, err)
			assert.Equal(t, tc.want, a.Compare(b))
			assert.Equal(t, -tc.want, b.Compare(a))
		})
	}
}

func TestRelations(t *testing.T) {
	a, _ := mvnversion.NewVersion("1.0")
	b, _ := mvnversion.NewVersion("2.0")

	assert.True(t, a.LessThan(b))
	assert.True(t, b.GreaterThan(a))
	assert.True(t, a.LessThanOrEqual(b))
	assert.True(t, b.GreaterThanOrEqual(a))
	assert.False(t, a.Equal(b))

	aeq, _ := mvnversion.NewVersion("1.0")
	assert.True(t, a.Equal(aeq))
	assert.True(t, a.GreaterThanOrEqual(aeq))
	assert.True(t, a.LessThanOrEqual(aeq))
}
