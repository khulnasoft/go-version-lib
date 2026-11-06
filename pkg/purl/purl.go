// Package purl implements Package URLs (purl) as defined by
// https://github.com/package-url/packageurl-spec
//
// A Package URL has the form:
//
//	pkg:type/namespace/name@version?qualifiers#subpath
//
// The package is intentionally independent of the version comparison
// packages so that it can be imported by them without import cycles.
package purl

import (
	"errors"
	"net/url"
	"sort"
	"strings"
)

const prefix = "pkg:"

// PackageURL is a parsed Package URL.
type PackageURL struct {
	// Type is the package type (e.g. "deb", "gem", "maven", "npm",
	// "pep440", "rpm", "semver").
	Type string

	// Namespace is the optional namespace. Multiple levels are joined
	// with "/".
	Namespace string

	// Name is the package name.
	Name string

	// Version is the optional package version.
	Version string

	// Qualifiers are optional key/value qualifiers.
	Qualifiers map[string]string

	// Subpath is the optional subpath.
	Subpath string
}

// Parse parses a Package URL string and returns a PackageURL.
func Parse(s string) (PackageURL, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, prefix) {
		return PackageURL{}, errors.New("missing pkg: prefix")
	}
	rest := strings.TrimPrefix(s, prefix)

	p := PackageURL{}

	// subpath: everything after the first '#'
	if i := strings.Index(rest, "#"); i >= 0 {
		p.Subpath = rest[i+1:]
		rest = rest[:i]
	}

	// qualifiers: everything after the first '?'
	if i := strings.Index(rest, "?"); i >= 0 {
		qs := rest[i+1:]
		rest = rest[:i]
		for _, q := range strings.Split(qs, "&") {
			if q == "" {
				continue
			}
			k, v, found := strings.Cut(q, "=")
			if !found {
				v = ""
			}
			dk, err := url.PathUnescape(k)
			if err != nil {
				return PackageURL{}, err
			}
			dv, err := url.PathUnescape(v)
			if err != nil {
				return PackageURL{}, err
			}
			if p.Qualifiers == nil {
				p.Qualifiers = map[string]string{}
			}
			p.Qualifiers[dk] = dv
		}
	}

	// version: the last segment is name[@version]; split on the first '@'.
	segments := strings.Split(rest, "/")
	if len(segments) == 0 || segments[0] == "" {
		return PackageURL{}, errors.New("missing package type")
	}
	p.Type = segments[0]

	if len(segments) == 1 {
		return PackageURL{}, errors.New("missing package name")
	}

	last := segments[len(segments)-1]
	if i := strings.Index(last, "@"); i >= 0 {
		p.Name = last[:i]
		p.Version = last[i+1:]
	} else {
		p.Name = last
	}
	if p.Name == "" {
		return PackageURL{}, errors.New("missing package name")
	}

	if len(segments) > 2 {
		p.Namespace = strings.Join(segments[1:len(segments)-1], "/")
	}

	var err error
	p.Namespace, err = url.PathUnescape(p.Namespace)
	if err != nil {
		return PackageURL{}, err
	}
	p.Name, err = url.PathUnescape(p.Name)
	if err != nil {
		return PackageURL{}, err
	}
	p.Version, err = url.PathUnescape(p.Version)
	if err != nil {
		return PackageURL{}, err
	}
	p.Subpath, err = url.PathUnescape(p.Subpath)
	if err != nil {
		return PackageURL{}, err
	}

	return p, nil
}

// encodeSegment percent-encodes a single purl segment. Characters that have
// structural meaning in a Package URL (: / @ ? # = &), percent (%), whitespace
// and control characters are encoded; everything else is left as-is so that
// common version/package strings (e.g. 1.0+build, 2!3.4, usr/bin) remain
// readable and round-trip byte-for-byte. Encoding the structural characters
// keeps the output unambiguous: Parse(encodeSegment(s)) == s.
func encodeSegment(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case ':', '/', '@', '?', '#', '=', '&', '%':
			sb.WriteString(url.PathEscape(string(r)))
		default:
			if r == ' ' || r < 0x20 || r == 0x7f {
				sb.WriteString(url.PathEscape(string(r)))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}

// String returns the canonical Package URL string representation of p.
// The output round-trips: Parse(p.String()) is equivalent to p.
func (p PackageURL) String() string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(p.Type)
	if p.Namespace != "" {
		for _, seg := range strings.Split(p.Namespace, "/") {
			sb.WriteString("/")
			sb.WriteString(encodeSegment(seg))
		}
	}
	sb.WriteString("/")
	sb.WriteString(encodeSegment(p.Name))
	if p.Version != "" {
		sb.WriteString("@")
		sb.WriteString(encodeSegment(p.Version))
	}
	if len(p.Qualifiers) > 0 {
		sb.WriteString("?")
		keys := make([]string, 0, len(p.Qualifiers))
		for k := range p.Qualifiers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if i > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(encodeSegment(k))
			sb.WriteString("=")
			sb.WriteString(encodeSegment(p.Qualifiers[k]))
		}
	}
	if p.Subpath != "" {
		sb.WriteString("#")
		segs := strings.Split(p.Subpath, "/")
		for i, seg := range segs {
			if i > 0 {
				sb.WriteString("/")
			}
			sb.WriteString(encodeSegment(seg))
		}
	}
	return sb.String()
}
