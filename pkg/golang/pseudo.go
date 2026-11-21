// Package golang provides helpers for Go module pseudo-versions.
package golang

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// pseudoVersion matches the three pseudo-version forms:
//
//	vX.0.0-yyyymmddhhmmss-abcdefabcdef
//	vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef
//	vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef
var pseudoVersion = regexp.MustCompile(`^v\d+\.\d+\.\d+-(?:\d{14}-[0-9a-f]{12}|(?:[0-9A-Za-z-]+\.)*0\.\d{14}-[0-9a-f]{12})(?:\+incompatible)?$`)

// timestamp extracts the commit timestamp from the pseudo-version tail.
var timestampRe = regexp.MustCompile(`(\d{14})-([0-9a-f]{12})(?:\+incompatible)?$`)

// IsPseudoVersion reports whether s is a Go module pseudo-version. The
// embedded timestamp must be a valid calendar time.
func IsPseudoVersion(s string) bool {
	if !pseudoVersion.MatchString(s) {
		return false
	}
	m := timestampRe.FindStringSubmatch(s)
	if m == nil {
		return false
	}
	_, err := time.Parse("20060102150405", m[1])
	return err == nil
}

// PseudoVersion is a parsed pseudo-version.
type PseudoVersion struct {
	// Base is the underlying semantic version portion (e.g. "v1.2.4-0" or
	// "v0.0.0").
	Base string
	// Timestamp is the commit time encoded in the version.
	Timestamp time.Time
	// Revision is the 12-character commit hash prefix.
	Revision string
}

// ParsePseudoVersion parses s and returns its components. It returns an
// error if s is not a pseudo-version.
func ParsePseudoVersion(s string) (PseudoVersion, error) {
	if !IsPseudoVersion(s) {
		return PseudoVersion{}, fmt.Errorf("not a pseudo-version: %s", s)
	}
	m := timestampRe.FindStringSubmatch(s)
	if m == nil {
		return PseudoVersion{}, fmt.Errorf("cannot extract timestamp from pseudo-version: %s", s)
	}
	ts, err := time.Parse("20060102150405", m[1])
	if err != nil {
		return PseudoVersion{}, fmt.Errorf("invalid pseudo-version timestamp: %w", err)
	}
	// Base is everything before the timestamp identifier's leading '.' or '-'
	idx := strings.LastIndex(s, m[1])
	base := strings.TrimRight(s[:idx], ".-")
	return PseudoVersion{Base: base, Timestamp: ts, Revision: m[2]}, nil
}
