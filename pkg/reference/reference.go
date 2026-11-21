package reference

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/khulnasoft/go-version-lib/pkg/golang"
	"github.com/khulnasoft/go-version-lib/pkg/purl"
)

type VersionType string

const (
	Semver   VersionType = "semver"
	Gem      VersionType = "gem"
	NPM      VersionType = "npm"
	PEP440   VersionType = "pep440"
	Deb      VersionType = "deb"
	RPM      VersionType = "rpm"
	Maven    VersionType = "maven"
	Cargo    VersionType = "cargo"
	Pub      VersionType = "pub"
	Hex      VersionType = "hex"
	NuGet    VersionType = "nuget"
	GoPseudo VersionType = "go-pseudo"
	Unknown  VersionType = "unknown"
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

var (
	registryMu  sync.Mutex
	registryGen uint64
	sortedCache []versionEntry
	sortedGen   uint64
)

func Register(typ VersionType, pattern string, parse parserFunc, order int) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registryGen++
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

// DetectAll returns every registered ecosystem whose pattern matches the
// input, in registration-priority order. Unlike Detect it does not stop at
// the first match, which is useful when version strings are ambiguous
// (e.g. "1.2.3" matches semver, npm, gem, and deb patterns).
func DetectAll(v string) []VersionType {
	trimmed := strings.TrimSpace(v)
	if strings.HasPrefix(trimmed, "pkg:") {
		if t, err := detectPurl(trimmed); err == nil {
			return []VersionType{t}
		}
		return nil
	}
	var out []VersionType
	for _, entry := range sortedEntries() {
		if entry.regex.MatchString(trimmed) {
			out = append(out, entry.typ)
		}
	}
	return out
}

func NewReferenceByType(v string, typ VersionType) (Reference, error) {
	registryMu.Lock()
	var parse parserFunc
	for _, entry := range registry {
		if entry.typ == typ {
			parse = entry.parse
			break
		}
	}
	registryMu.Unlock()
	if parse == nil {
		return nil, fmt.Errorf("unknown version type: %s", typ)
	}
	// invoke the parser outside the lock: a registered parser may itself
	// call back into the registry
	return parse(v)
}

// sortedEntries returns the registry sorted by registration order. The sorted
// view is cached and only rebuilt when the registry changes.
func sortedEntries() []versionEntry {
	registryMu.Lock()
	defer registryMu.Unlock()
	if sortedGen != registryGen {
		sortedCache = make([]versionEntry, len(registry))
		copy(sortedCache, registry)
		sort.SliceStable(sortedCache, func(i, j int) bool {
			return sortedCache[i].order < sortedCache[j].order
		})
		sortedGen = registryGen
	}
	return sortedCache
}

// compareDifferentTypes provides a deterministic total order when comparing
// references from different ecosystems (or mixed purl/ecosystem references).
func compareDifferentTypes(a, b Reference) int {
	if c := strings.Compare(string(a.Type()), string(b.Type())); c != 0 {
		return c
	}
	return strings.Compare(a.String(), b.String())
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
	case "cargo":
		return Cargo
	case "pub":
		return Pub
	case "hex":
		return Hex
	case "nuget":
		return NuGet
	case "golang":
		return GoPseudo
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
	if t == GoPseudo && !golang.IsPseudoVersion(p.Version) {
		// ordinary tagged Go releases are not pseudo-versions; Go module
		// versions carry a "v" prefix that semver.Parse rejects
		t = Semver
		p.Version = strings.TrimPrefix(p.Version, "v")
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
		return compareDifferentTypes(r, other)
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
