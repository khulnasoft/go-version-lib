package reference

import (
	"strings"

	cargoversion "github.com/khulnasoft/go-version-lib/pkg/cargo"
	debversion "github.com/khulnasoft/go-version-lib/pkg/deb"
	"github.com/khulnasoft/go-version-lib/pkg/gem"
	"github.com/khulnasoft/go-version-lib/pkg/golang"
	hexversion "github.com/khulnasoft/go-version-lib/pkg/hex"
	mvnversion "github.com/khulnasoft/go-version-lib/pkg/mvn"
	"github.com/khulnasoft/go-version-lib/pkg/npm"
	"github.com/khulnasoft/go-version-lib/pkg/nuget"
	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
	pubversion "github.com/khulnasoft/go-version-lib/pkg/pub"
	rpmversion "github.com/khulnasoft/go-version-lib/pkg/rpm"
	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

const (
	semverOrder   = 100
	gemOrder      = 200
	npmOrder      = 300
	pep440Order   = 400
	debOrder      = 500
	rpmOrder      = 600
	mvnOrder      = 700
	cargoOrder    = 800
	pubOrder      = 900
	hexOrder      = 1000
	nugetOrder    = 1100
	goPseudoOrder = 50
)

func init() {
	Register(Semver, semverRegex, func(v string) (Reference, error) {
		ver, err := semver.Parse(v)
		if err != nil {
			return nil, err
		}
		return &semverRef{ver}, nil
	}, semverOrder)

	Register(Gem, gemRegex, func(v string) (Reference, error) {
		ver, err := gem.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &gemRef{ver}, nil
	}, gemOrder)

	Register(NPM, npmRegex, func(v string) (Reference, error) {
		ver, err := npm.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &npmRef{ver}, nil
	}, npmOrder)

	Register(PEP440, pep440Regex, func(v string) (Reference, error) {
		ver, err := pep440version.Parse(v)
		if err != nil {
			return nil, err
		}
		return &pep440Ref{ver}, nil
	}, pep440Order)

	Register(Deb, debRegex, func(v string) (Reference, error) {
		ver, err := debversion.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &debRef{ver}, nil
	}, debOrder)

	Register(RPM, rpmRegex, func(v string) (Reference, error) {
		ver := rpmversion.NewVersion(v)
		return &rpmRef{ver}, nil
	}, rpmOrder)

	Register(Maven, mavenRegex, func(v string) (Reference, error) {
		ver, err := mvnversion.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &mvnRef{ver}, nil
	}, mvnOrder)

	Register(Cargo, semverRegex, func(v string) (Reference, error) {
		ver, err := cargoversion.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &cargoRef{ver}, nil
	}, cargoOrder)

	Register(Pub, semverRegex, func(v string) (Reference, error) {
		ver, err := pubversion.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &pubRef{ver}, nil
	}, pubOrder)

	Register(Hex, semverRegex, func(v string) (Reference, error) {
		ver, err := hexversion.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &hexRef{ver}, nil
	}, hexOrder)

	Register(NuGet, nugetRegex, func(v string) (Reference, error) {
		ver, err := nuget.NewVersion(v)
		if err != nil {
			return nil, err
		}
		return &nugetRef{ver}, nil
	}, nugetOrder)

	Register(GoPseudo, goPseudoRegex, func(v string) (Reference, error) {
		if _, err := golang.ParsePseudoVersion(v); err != nil {
			return nil, err
		}
		ver, err := semver.Parse(strings.TrimPrefix(v, "v"))
		if err != nil {
			return nil, err
		}
		return &goPseudoRef{ver: ver, original: v}, nil
	}, goPseudoOrder)
}

var semverRegex = `^v?([0-9]+)\.([0-9]+)\.([0-9]+)` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

var gemRegex = `^[0-9]+(\.[0-9a-zA-Z]+)*(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`

var npmRegex = `^v?([0-9|x|X|\*]+)(\.[0-9|x|X|\*]+)?(\.[0-9|x|X|\*]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

var pep440Regex = `^v?` +
	`(?:[0-9]+!)?` +
	`[0-9]+(?:\.[0-9]+)*` +
	`(?:[-_\.]?(?:a|b|c|rc|alpha|beta|pre|preview)[-_\.]?[0-9]?)?` +
	`(?:-(?:[0-9]+)|[-_\.]?(?:post|rev|r)[-_\.]?[0-9]?)?` +
	`(?:[-_\.]?(?:dev)[-_\.]?[0-9]?)?` +
	`(?:\+[a-z0-9]+(?:[-_\.][a-z0-9]+)?)?$`

var debRegex = `^([0-9]+:)?[0-9]+(\.[0-9a-zA-Z]+)*(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`

var rpmRegex = `^([0-9]+:)?[0-9A-Za-z\.\+\~\^]+(-[0-9A-Za-z\.\+\~\^]+)?$`

var mavenRegex = `^[0-9A-Za-z\-~_\.\+]+$`

var nugetRegex = `^[0-9]+(\.[0-9]+){0,3}(-[0-9A-Za-z.-]+)?$`

var goPseudoRegex = `^v[0-9]+\.[0-9]+\.[0-9]+-(?:[0-9]{14}-[0-9a-f]{12}|(?:[0-9A-Za-z-]+\.)*0\.[0-9]{14}-[0-9a-f]{12})(?:\+incompatible)?$`

type semverRef struct {
	ver semver.Version
}

func (r *semverRef) String() string { return r.ver.String() }

func (r *semverRef) Compare(other Reference) int {
	o, ok := other.(*semverRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *semverRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *semverRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *semverRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *semverRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *semverRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *semverRef) Type() VersionType                       { return Semver }

type gemRef struct {
	ver gem.Version
}

func (r *gemRef) String() string { return r.ver.String() }

func (r *gemRef) Compare(other Reference) int {
	o, ok := other.(*gemRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *gemRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *gemRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *gemRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *gemRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *gemRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *gemRef) Type() VersionType                       { return Gem }

type npmRef struct {
	ver npm.Version
}

func (r *npmRef) String() string { return r.ver.String() }

func (r *npmRef) Compare(other Reference) int {
	o, ok := other.(*npmRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *npmRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *npmRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *npmRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *npmRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *npmRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *npmRef) Type() VersionType                       { return NPM }

type pep440Ref struct {
	ver pep440version.Version
}

func (r *pep440Ref) String() string { return r.ver.String() }

func (r *pep440Ref) Compare(other Reference) int {
	o, ok := other.(*pep440Ref)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *pep440Ref) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *pep440Ref) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *pep440Ref) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *pep440Ref) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *pep440Ref) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *pep440Ref) Type() VersionType                       { return PEP440 }

type debRef struct {
	ver debversion.Version
}

func (r *debRef) String() string { return r.ver.String() }

func (r *debRef) Compare(other Reference) int {
	o, ok := other.(*debRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *debRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *debRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *debRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *debRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *debRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *debRef) Type() VersionType                       { return Deb }

type rpmRef struct {
	ver rpmversion.Version
}

func (r *rpmRef) String() string { return r.ver.String() }

func (r *rpmRef) Compare(other Reference) int {
	o, ok := other.(*rpmRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *rpmRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *rpmRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *rpmRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *rpmRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *rpmRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *rpmRef) Type() VersionType                       { return RPM }

type mvnRef struct {
	ver mvnversion.Version
}

func (r *mvnRef) String() string { return r.ver.String() }

func (r *mvnRef) Compare(other Reference) int {
	o, ok := other.(*mvnRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *mvnRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *mvnRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *mvnRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *mvnRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *mvnRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *mvnRef) Type() VersionType                       { return Maven }

type cargoRef struct {
	ver cargoversion.Version
}

func (r *cargoRef) String() string { return r.ver.String() }

func (r *cargoRef) Compare(other Reference) int {
	o, ok := other.(*cargoRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *cargoRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *cargoRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *cargoRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *cargoRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *cargoRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *cargoRef) Type() VersionType                       { return Cargo }

type pubRef struct {
	ver pubversion.Version
}

func (r *pubRef) String() string { return r.ver.String() }

func (r *pubRef) Compare(other Reference) int {
	o, ok := other.(*pubRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *pubRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *pubRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *pubRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *pubRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *pubRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *pubRef) Type() VersionType                       { return Pub }

type hexRef struct {
	ver hexversion.Version
}

func (r *hexRef) String() string { return r.ver.String() }

func (r *hexRef) Compare(other Reference) int {
	o, ok := other.(*hexRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *hexRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *hexRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *hexRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *hexRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *hexRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *hexRef) Type() VersionType                       { return Hex }

type nugetRef struct {
	ver nuget.Version
}

func (r *nugetRef) String() string { return r.ver.String() }

func (r *nugetRef) Compare(other Reference) int {
	o, ok := other.(*nugetRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *nugetRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *nugetRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *nugetRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *nugetRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *nugetRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *nugetRef) Type() VersionType                       { return NuGet }

type goPseudoRef struct {
	ver      semver.Version
	original string
}

func (r *goPseudoRef) String() string { return r.original }

func (r *goPseudoRef) Compare(other Reference) int {
	o, ok := other.(*goPseudoRef)
	if !ok {
		return compareDifferentTypes(r, other)
	}
	return r.ver.Compare(o.ver)
}

func (r *goPseudoRef) Equal(other Reference) bool              { return r.Compare(other) == 0 }
func (r *goPseudoRef) GreaterThan(other Reference) bool        { return r.Compare(other) > 0 }
func (r *goPseudoRef) LessThan(other Reference) bool           { return r.Compare(other) < 0 }
func (r *goPseudoRef) GreaterThanOrEqual(other Reference) bool { return r.Compare(other) >= 0 }
func (r *goPseudoRef) LessThanOrEqual(other Reference) bool    { return r.Compare(other) <= 0 }
func (r *goPseudoRef) Type() VersionType                       { return GoPseudo }
