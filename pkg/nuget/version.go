// Package nuget implements NuGet package version parsing and comparison.
// NuGet versions are up to four numeric release segments
// (Major.Minor.Build.Revision) with an optional prerelease label.
package nuget

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var versionRegex = regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z.-]+))?$`)

// Version is a NuGet package version.
type Version struct {
	parts      [4]uint64
	prerelease string
	metadata   string
	original   string
}

// NewVersion parses a NuGet version string such as "1.2.3.4-beta.2".
func NewVersion(s string) (Version, error) {
	m := versionRegex.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Version{}, fmt.Errorf("malformed nuget version: %s", s)
	}
	var v Version
	for i := 0; i < 4; i++ {
		if m[i+1] == "" {
			continue
		}
		n, err := strconv.ParseUint(m[i+1], 10, 64)
		if err != nil {
			return Version{}, fmt.Errorf("malformed nuget version %q: %w", s, err)
		}
		v.parts[i] = n
	}
	if pre := m[5]; pre != "" {
		for _, ident := range strings.Split(pre, ".") {
			if isNumericIdent(ident) && len(ident) > 1 && ident[0] == '0' {
				return Version{}, fmt.Errorf("malformed nuget version %q: numeric prerelease identifier %q has leading zeros", s, ident)
			}
		}
	}
	v.prerelease = m[5]
	v.metadata = m[6]
	v.original = s
	return v, nil
}

// MustParse parses s and panics on error.
func MustParse(s string) Version {
	v, err := NewVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns the original version string.
func (v Version) String() string { return v.original }

// Prerelease returns the prerelease label, or "" for a release version.
func (v Version) Prerelease() string { return v.prerelease }

// Metadata returns the build metadata suffix, or "" if none. Build
// metadata is preserved by String but ignored by Compare, matching NuGet
// precedence rules.
func (v Version) Metadata() string { return v.metadata }

// Compare orders versions numerically by release segment, then by
// prerelease label (a release sorts after any of its prereleases).
func (v Version) Compare(o Version) int {
	for i := 0; i < 4; i++ {
		if v.parts[i] != o.parts[i] {
			if v.parts[i] < o.parts[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case v.prerelease == "" && o.prerelease == "":
		return 0
	case v.prerelease == "":
		return 1
	case o.prerelease == "":
		return -1
	}
	return comparePrerelease(v.prerelease, o.prerelease)
}

// comparePrerelease orders dot-separated identifiers: numeric identifiers
// compare numerically and sort before alphanumeric ones; alphanumeric
// identifiers compare lexically (case-insensitive); a shorter prefix sorts
// first.
func isNumericIdent(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func comparePrerelease(a, b string) int {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		if isNumericIdent(as[i]) && isNumericIdent(bs[i]) {
			// compare by numeric value without a uint64 limit: longer
			// digit strings (after trimming zeros) are larger
			an, bn := strings.TrimLeft(as[i], "0"), strings.TrimLeft(bs[i], "0")
			if len(an) != len(bn) {
				if len(an) < len(bn) {
					return -1
				}
				return 1
			}
			if an != bn {
				return strings.Compare(an, bn)
			}
			continue
		}
		_, aerr := strconv.ParseUint(as[i], 10, 64)
		_, berr := strconv.ParseUint(bs[i], 10, 64)
		switch {
		case isNumericIdent(as[i]):
			return -1
		case isNumericIdent(bs[i]):
			return 1
		case aerr == nil:
			return -1
		case berr == nil:
			return 1
		default:
			if c := strings.Compare(strings.ToLower(as[i]), strings.ToLower(bs[i])); c != 0 {
				return c
			}
		}
	}
	return len(as) - len(bs)
}

func (v Version) Equal(o Version) bool              { return v.Compare(o) == 0 }
func (v Version) GreaterThan(o Version) bool        { return v.Compare(o) > 0 }
func (v Version) LessThan(o Version) bool           { return v.Compare(o) < 0 }
func (v Version) GreaterThanOrEqual(o Version) bool { return v.Compare(o) >= 0 }
func (v Version) LessThanOrEqual(o Version) bool    { return v.Compare(o) <= 0 }
