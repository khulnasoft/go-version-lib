package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

// coerceRegex extracts the first numeric version-looking substring,
// tolerating a leading "v" and missing minor/patch segments.
var coerceRegex = regexp.MustCompile(`v?(\d+)(?:\.(\d+))?(?:\.(\d+))?`)

// Coerce extracts a semantic version from an arbitrary string, padding
// missing minor/patch segments with zero. It is more permissive than Parse:
// "v1.2" coerces to "1.2.0" and "foo-1.2.3-bar" to "1.2.3".
func Coerce(s string) (Version, error) {
	m := coerceRegex.FindStringSubmatch(s)
	if m == nil {
		return Version{}, fmt.Errorf("cannot coerce version: %q", s)
	}
	major, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return Version{}, fmt.Errorf("cannot coerce version %q: %w", s, err)
	}
	var minor, patch uint64
	if m[2] != "" {
		minor, err = strconv.ParseUint(m[2], 10, 64)
		if err != nil {
			return Version{}, fmt.Errorf("cannot coerce version %q: %w", s, err)
		}
	}
	if m[3] != "" {
		patch, err = strconv.ParseUint(m[3], 10, 64)
		if err != nil {
			return Version{}, fmt.Errorf("cannot coerce version %q: %w", s, err)
		}
	}
	return Parse(fmt.Sprintf("%d.%d.%d", major, minor, patch))
}

// MustCoerce is like Coerce but panics on error.
func MustCoerce(s string) Version {
	v, err := Coerce(s)
	if err != nil {
		panic(err)
	}
	return v
}

// CoerceOrNil returns the coerced version and whether coercion succeeded.
func CoerceOrNil(s string) (Version, bool) {
	v, err := Coerce(s)
	if err != nil {
		return Version{}, false
	}
	return v, true
}
