package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	versionlib "github.com/khulnasoft/go-version-lib/pkg/version"
)

func TestVersionOriginal(t *testing.T) {
	v, err := versionlib.Parse("v1.2.3-beta+build")
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3-beta+build", v.Original())
	assert.Equal(t, "1.2.3-beta+build", v.String())
}

func TestVersionConstraintStringAndOperators(t *testing.T) {
	c, err := versionlib.NewConstraints(">= 1.0.0, < 2.0.0")
	require.NoError(t, err)

	s := c.String()
	assert.Contains(t, s, ">= 1.0.0")
	assert.Contains(t, s, "< 2.0.0")

	v, err := versionlib.Parse("1.5.0")
	require.NoError(t, err)
	assert.True(t, c.Check(v))

	v2, err := versionlib.Parse("2.0.0")
	require.NoError(t, err)
	assert.False(t, c.Check(v2))
}

func TestVersionConstraintSingleOperator(t *testing.T) {
	// Test each operator
	ops := []struct {
		op   string
		want bool
	}{
		{"= 1.0.0", true},
		{"== 1.0.0", true},
		{"!= 1.0.0", false},
		{"> 1.0.0", false},
		{"< 1.0.0", false},
		{">= 1.0.0", true},
		{"<= 1.0.0", true},
		{"=> 1.0.0", true},
		{"=< 1.0.0", true},
	}

	v, _ := versionlib.Parse("1.0.0")
	for _, tc := range ops {
		t.Run(tc.op, func(t *testing.T) {
			c, err := versionlib.NewConstraints(tc.op)
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Check(v))
		})
	}
}

func TestVersionConstraintExportedMethods(t *testing.T) {
	c, err := versionlib.NewConstraints(">= 1.0.0, < 2.0.0")
	require.NoError(t, err)

	// List
	lists := c.List()
	assert.Len(t, lists, 1)
	assert.Len(t, lists[0], 2)

	// Constraint methods
	constraint := lists[0][0]
	assert.Equal(t, ">=", constraint.Operator())
	assert.Equal(t, "1.0.0", constraint.Version())
	assert.Equal(t, ">= 1.0.0", constraint.String())
}
