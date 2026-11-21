// Package hex provides semantic version support for the hex ecosystem.
// Versions and constraints are aliases of the semver implementation since
// these ecosystems use Semantic Versioning.
package hex

import (
	"github.com/khulnasoft/go-version-lib/pkg/semver"
	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// Version represents a semantic version.
type Version = semver.Version

// NewVersion parses a given version and returns an instance of Version.
func NewVersion(s string) (Version, error) {
	return semver.Parse(s)
}

// Coerce extracts a version from an arbitrary string.
func Coerce(s string) (Version, error) { return semver.Coerce(s) }

// Constraints is a set of semver constraints.
type Constraints = semver.Constraints

// ConstraintOption configures constraint parsing.
type ConstraintOption = semver.ConstraintOption

// NewConstraints parses a Hex (Mix) requirement expression. It supports
// "and"/"or" keywords and "~>" pessimistic requirements, translated to
// semver constraint syntax before parsing.
func NewConstraints(v string, opts ...ConstraintOption) (Constraints, error) {
	translated, err := translateHexConstraint(v)
	if err != nil {
		return Constraints{}, err
	}
	return semver.NewConstraints(translated, opts...)
}

// ToVersionRange lowers constraints to the generic versionrange model.
func ToVersionRange(cs Constraints) (versionrange.Range[Version], error) {
	return cs.ToVersionRange()
}
