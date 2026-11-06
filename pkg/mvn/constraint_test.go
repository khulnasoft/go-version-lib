package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mvnversion "github.com/khulnasoft/go-version-lib/pkg/mvn"
)

func TestMvnNewConstraints(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
		err        bool
	}{
		{"1.0", "1.0", true, false},
		{"1.0", "2.0", false, false},
		{">= 1.0", "1.0", true, false},
		{">= 1.0", "0.9", false, false},
		{"< 2.0", "1.9", true, false},
		{"< 2.0", "2.0", false, false},
		{"improper [", "", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.constraint+"|"+tc.version, func(t *testing.T) {
			c, err := mvnversion.NewConstraints(tc.constraint)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			v, err := mvnversion.NewVersion(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Check(v))
		})
	}
}

func TestMvnConstraintCheck(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		{"1.0", "1.0", true},
		{"1.0", "2.0", false},
		{">= 1.0", "1.0", true},
		{">= 1.0", "0.9", false},
		{"< 2.0", "1.9", true},
		{"< 2.0", "2.0", false},
		{"!= 1.0", "1.1", true},
		{"!= 1.0", "1.0", false},
	}

	for _, tc := range tests {
		t.Run(tc.constraint+"|"+tc.version, func(t *testing.T) {
			c, err := mvnversion.NewConstraints(tc.constraint)
			require.NoError(t, err)
			v, err := mvnversion.NewVersion(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Check(v))
		})
	}
}

func TestMvnConstraintString(t *testing.T) {
	c, err := mvnversion.NewConstraints("1.0, < 2.0")
	require.NoError(t, err)
	assert.Contains(t, c.String(), "1.0")
	assert.Contains(t, c.String(), "< 2.0")
}

func TestMvnNewRequirements(t *testing.T) {
	tests := []struct {
		in  string
		err bool
	}{
		{"[1.0.0]", false},
		{"(1.0.0)", false},
		{"[1.0.0, 2.0.0]", false},
		{"[,1.0.0]", false},
		{"[1.0.0,)", false},
		{"1.0.0", false},
		{"[1.0.0],[1.0.1]", false},
		{"improper [", true},
		{"[1.0.0, 2.0.0, 3.0.0]", true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			_, err := mvnversion.NewRequirements(tc.in)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMvnRequirementCheck(t *testing.T) {
	tests := []struct {
		requirement string
		version     string
		want        bool
	}{
		{"[1.0.0]", "1.0.0", true},
		{"[1.0.0]", "2.0.0", false},
		{"(1.0.0)", "1.0.0", true},
		{"(1.0.0)", "2.0.0", false},
		{"[1.0.0, 2.0.0]", "1.5.0", true},
		{"[1.0.0, 2.0.0]", "3.0.0", false},
		{"[,1.0.0]", "0.9.0", true},
		{"[,1.0.0]", "2.0.0", false},
		{"[1.0.0,)", "2.0.0", true},
		{"[1.0.0,)", "0.9.0", false},
	}

	for _, tc := range tests {
		t.Run(tc.requirement+"|"+tc.version, func(t *testing.T) {
			r, err := mvnversion.NewRequirements(tc.requirement)
			require.NoError(t, err)
			v, err := mvnversion.NewVersion(tc.version)
			require.NoError(t, err)
			assert.Equal(t, tc.want, r.Check(v))
		})
	}
}

func TestMvnNewComparer(t *testing.T) {
	// constraint syntax
	c, err := mvnversion.NewComparer(">= 1.0")
	require.NoError(t, err)
	v, _ := mvnversion.NewVersion("2.0")
	assert.True(t, c.Check(v))

	// requirement syntax
	c, err = mvnversion.NewComparer("[1.0.0]")
	require.NoError(t, err)
	v, _ = mvnversion.NewVersion("1.0.0")
	assert.True(t, c.Check(v))

	// invalid
	_, err = mvnversion.NewComparer("invalid [[[")
	require.Error(t, err)
}

func TestMvnCollection(t *testing.T) {
	v1, _ := mvnversion.NewVersion("1.0.0")
	v2, _ := mvnversion.NewVersion("2.0.0")

	// Use a slice directly since mvn has no Collection type.
	versions := []mvnversion.Version{v2, v1}
	assert.Equal(t, 2, len(versions))
	assert.True(t, versions[1].LessThan(versions[0]))
}
