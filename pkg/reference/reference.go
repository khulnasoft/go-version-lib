package reference

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/khulnasoft/go-version-lib/pkg/purl"
)

type VersionType string

const (
	Semver  VersionType = "semver"
	Gem     VersionType = "gem"
	NPM     VersionType = "npm"
	PEP440  VersionType = "pep440"
	Deb     VersionType = "deb"
	RPM     VersionType = "rpm"
	Maven   VersionType = "maven"
	Unknown VersionType = "unknown"
)

type Reference interface {
	String() string
	Compare(other Reference) int
	Equal(other Reference) bool
	GreaterThan(other Reference) bool
	LessThan(other Reference) bool
	GreaterThanOrEqual(other Reference) bool
	LessThanOrEqual(other Reference) bool
	Type() VersionType
}

type parserFunc func(string) (Reference, error)

type versionEntry struct {
	typ   VersionType
	regex *regexp.Regexp
	parse parserFunc
	order int
}

var registry []versionEntry

func Register(typ VersionType, pattern string, parse parserFunc, order int) {
	for i, entry := range registry {
		if entry.typ == typ {
			registry[i] = versionEntry{
				typ:   typ,
				regex: regexp.MustCompile(pattern),
				parse: parse,
				order: order,
			}
			return
		}
	}
	registry = append(registry, versionEntry{
		typ:   typ,
		regex: regexp.MustCompile(pattern),
		parse: parse,
		order: order,
	})
}

func Detect(v string) (VersionType, error) {
	trimmed := strings.TrimSpace(v)
	if strings.HasPrefix(trimmed, "pkg:") {
		return detectPurl(trimmed)
	}
	for _, entry := range sortedEntries() {
		if entry.regex.MatchString(trimmed) {
			return entry.typ, nil
		}
	}
	return Unknown, fmt.Errorf("cannot detect version type for: %s", v)
}

func NewReference(v string) (Reference, error) {
	trimmed := strings.TrimSpace(v)
	if strings.HasPrefix(trimmed, "pkg:") {
		return newPurlReference(trimmed)
	}
	for _, entry := range sortedEntries() {
		if entry.regex.MatchString(trimmed) {
			return entry.parse(trimmed)
		}
	}
	return nil, fmt.Errorf("no parser matched for version: %s", v)
}

func NewReferenceByType(v string, typ VersionType) (Reference, error) {
	for _, entry := range registry {
		if entry.typ == typ {
			return entry.parse(v)
		}
	}
	return nil, fmt.Errorf("unknown version type: %s", typ)
}

func sortedEntries() []versionEntry {
	entries := make([]versionEntry, len(registry))
	copy(entries, registry)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].order > entries[j].order {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	return entries
}

// purlVersionType maps a purl package type to a reference VersionType.
func purlVersionType(t string) VersionType {
	switch t {
	case "deb":
		return Deb
	case "gem":
		return Gem
	case "maven":
		return Maven
	case "npm":
		return NPM
	case "pep440":
		return PEP440
	case "rpm":
		return RPM
	case "semver":
		return Semver
	default:
		return Unknown
	}
}

func detectPurl(s string) (VersionType, error) {
	p, err := purl.Parse(s)
	if err != nil {
		return Unknown, err
	}
	t := purlVersionType(p.Type)
	if t == Unknown {
		return Unknown, fmt.Errorf("unknown package type: %s", p.Type)
	}
	return t, nil
}

func newPurlReference(s string) (Reference, error) {
	p, err := purl.Parse(s)
	if err != nil {
		return nil, err
	}
	if p.Version == "" {
		return nil, fmt.Errorf("package url has no version: %s", s)
	}
	t := purlVersionType(p.Type)
	if t == Unknown {
		return nil, fmt.Errorf("unknown package type: %s", p.Type)
	}
	ref, err := NewReferenceByType(p.Version, t)
	if err != nil {
		return nil, err
	}
	return &purlRef{purl: p, ref: ref}, nil
}

// purlRef wraps an ecosystem Reference parsed from a Package URL. String()
// round-trips the canonical Package URL so the package identity is preserved.
type purlRef struct {
	purl purl.PackageURL
	ref  Reference
}

func (r *purlRef) String() string                          { return r.purl.String() }
func (r *purlRef) Compare(other Reference) int             { return comparePurlRef(r, other) }
func (r *purlRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *purlRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *purlRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *purlRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *purlRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *purlRef) Type() VersionType                       { return purlVersionType(r.purl.Type) }

// Purl returns the underlying Package URL.
func (r *purlRef) Purl() purl.PackageURL { return r.purl }

// Version returns the underlying parsed ecosystem version reference.
func (r *purlRef) Version() Reference { return r.ref }

func comparePurlRef(r *purlRef, other Reference) int {
	o, ok := other.(*purlRef)
	if !ok {
		return strings.Compare(r.String(), other.String())
	}
	// Same package identity: compare the ecosystem versions version-aware.
	if r.purl.Type == o.purl.Type &&
		r.purl.Namespace == o.purl.Namespace &&
		r.purl.Name == o.purl.Name {
		return r.ref.Compare(o.ref)
	}
	// Different packages: order by package identity, then version.
	if c := strings.Compare(r.purl.Type, o.purl.Type); c != 0 {
		return c
	}
	if c := strings.Compare(r.purl.Namespace, o.purl.Namespace); c != 0 {
		return c
	}
	if c := strings.Compare(r.purl.Name, o.purl.Name); c != 0 {
		return c
	}
	return r.ref.Compare(o.ref)
}
