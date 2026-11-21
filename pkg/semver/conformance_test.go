package semver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

// TestSemverOrgExamples exercises the version examples published at
// https://semver.org — both the valid and invalid lists.
func TestSemverOrgExamples(t *testing.T) {
	valid := []string{
		"0.0.4",
		"1.2.3",
		"10.20.30",
		"1.1.2-prerelease+meta",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
		"2.0.0",
		"2.1.1",
		"5.3.21",
		"1.0.0-0.3.7",
		"1.0.0-x.7.z.92",
		"1.0.0-x-y-z.--",
	}
	for _, s := range valid {
		t.Run("valid/"+s, func(t *testing.T) {
			_, err := semver.Parse(s)
			require.NoError(t, err)
		})
	}

	invalid := []string{
		"1",
		"1.2",
		"1.2.3.4",
		"01.2.3",
		"1.02.3",
		"1.2.03",
		"1.2.3-",
		"1.2.3-alpha..",
		"1.2.3-alpha_beta",
		"1.2.3-alpha+",
		"1.2.3+meta..x",
		"1.2.3-01",
		"v1.2.3",
		" 1.2.3",
		"1.2.3 ",
	}
	for _, s := range invalid {
		t.Run("invalid/"+s, func(t *testing.T) {
			_, err := semver.Parse(s)
			require.Error(t, err)
		})
	}
}

// TestSemverOrgPrecedence covers the complete precedence chain from the
// semver.org specification (item 11).
func TestSemverOrgPrecedence(t *testing.T) {
	chain := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
	}
	versions := make([]semver.Version, len(chain))
	for i, s := range chain {
		v, err := semver.Parse(s)
		require.NoError(t, err)
		versions[i] = v
	}
	for i := 0; i+1 < len(versions); i++ {
		assert.Negative(t, versions[i].Compare(versions[i+1]), "%s should be < %s", chain[i], chain[i+1])
		assert.Positive(t, versions[i+1].Compare(versions[i]), "%s should be > %s", chain[i+1], chain[i])
	}
}
