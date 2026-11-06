package reference

import (
	debversion "github.com/khulnasoft/go-version-lib/pkg/deb"
	"github.com/khulnasoft/go-version-lib/pkg/gem"
	mvnversion "github.com/khulnasoft/go-version-lib/pkg/mvn"
	"github.com/khulnasoft/go-version-lib/pkg/npm"
	pep440version "github.com/khulnasoft/go-version-lib/pkg/pep440"
	rpmversion "github.com/khulnasoft/go-version-lib/pkg/rpm"
	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

const (
	semverOrder = 100
	gemOrder    = 200
	npmOrder    = 300
	pep440Order = 400
	debOrder    = 500
	rpmOrder    = 600
	mvnOrder    = 700
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

type semverRef struct {
	ver semver.Version
}

func (r *semverRef) String() string              { return r.ver.String() }
func (r *semverRef) Compare(other Reference) int { return r.ver.Compare(other.(*semverRef).ver) }
func (r *semverRef) Equal(other Reference) bool  { return r.ver.Equal(other.(*semverRef).ver) }
func (r *semverRef) GreaterThan(other Reference) bool {
	return r.ver.GreaterThan(other.(*semverRef).ver)
}
func (r *semverRef) LessThan(other Reference) bool { return r.ver.LessThan(other.(*semverRef).ver) }
func (r *semverRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*semverRef).ver)
}
func (r *semverRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*semverRef).ver)
}
func (r *semverRef) Type() VersionType { return Semver }

type gemRef struct {
	ver gem.Version
}

func (r *gemRef) String() string                   { return r.ver.String() }
func (r *gemRef) Compare(other Reference) int      { return r.ver.Compare(other.(*gemRef).ver) }
func (r *gemRef) Equal(other Reference) bool       { return r.ver.Equal(other.(*gemRef).ver) }
func (r *gemRef) GreaterThan(other Reference) bool { return r.ver.GreaterThan(other.(*gemRef).ver) }
func (r *gemRef) LessThan(other Reference) bool    { return r.ver.LessThan(other.(*gemRef).ver) }
func (r *gemRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*gemRef).ver)
}
func (r *gemRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*gemRef).ver)
}
func (r *gemRef) Type() VersionType { return Gem }

type npmRef struct {
	ver npm.Version
}

func (r *npmRef) String() string                   { return r.ver.String() }
func (r *npmRef) Compare(other Reference) int      { return r.ver.Compare(other.(*npmRef).ver) }
func (r *npmRef) Equal(other Reference) bool       { return r.ver.Equal(other.(*npmRef).ver) }
func (r *npmRef) GreaterThan(other Reference) bool { return r.ver.GreaterThan(other.(*npmRef).ver) }
func (r *npmRef) LessThan(other Reference) bool    { return r.ver.LessThan(other.(*npmRef).ver) }
func (r *npmRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*npmRef).ver)
}
func (r *npmRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*npmRef).ver)
}
func (r *npmRef) Type() VersionType { return NPM }

type pep440Ref struct {
	ver pep440version.Version
}

func (r *pep440Ref) String() string              { return r.ver.String() }
func (r *pep440Ref) Compare(other Reference) int { return r.ver.Compare(other.(*pep440Ref).ver) }
func (r *pep440Ref) Equal(other Reference) bool  { return r.ver.Equal(other.(*pep440Ref).ver) }
func (r *pep440Ref) GreaterThan(other Reference) bool {
	return r.ver.GreaterThan(other.(*pep440Ref).ver)
}
func (r *pep440Ref) LessThan(other Reference) bool { return r.ver.LessThan(other.(*pep440Ref).ver) }
func (r *pep440Ref) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*pep440Ref).ver)
}
func (r *pep440Ref) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*pep440Ref).ver)
}
func (r *pep440Ref) Type() VersionType { return PEP440 }

type debRef struct {
	ver debversion.Version
}

func (r *debRef) String() string                   { return r.ver.String() }
func (r *debRef) Compare(other Reference) int      { return r.ver.Compare(other.(*debRef).ver) }
func (r *debRef) Equal(other Reference) bool       { return r.ver.Equal(other.(*debRef).ver) }
func (r *debRef) GreaterThan(other Reference) bool { return r.ver.GreaterThan(other.(*debRef).ver) }
func (r *debRef) LessThan(other Reference) bool    { return r.ver.LessThan(other.(*debRef).ver) }
func (r *debRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*debRef).ver)
}
func (r *debRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*debRef).ver)
}
func (r *debRef) Type() VersionType { return Deb }

type rpmRef struct {
	ver rpmversion.Version
}

func (r *rpmRef) String() string                   { return r.ver.String() }
func (r *rpmRef) Compare(other Reference) int      { return r.ver.Compare(other.(*rpmRef).ver) }
func (r *rpmRef) Equal(other Reference) bool       { return r.ver.Equal(other.(*rpmRef).ver) }
func (r *rpmRef) GreaterThan(other Reference) bool { return r.ver.GreaterThan(other.(*rpmRef).ver) }
func (r *rpmRef) LessThan(other Reference) bool    { return r.ver.LessThan(other.(*rpmRef).ver) }
func (r *rpmRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*rpmRef).ver)
}
func (r *rpmRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*rpmRef).ver)
}
func (r *rpmRef) Type() VersionType { return RPM }

type mvnRef struct {
	ver mvnversion.Version
}

func (r *mvnRef) String() string                   { return r.ver.String() }
func (r *mvnRef) Compare(other Reference) int      { return r.ver.Compare(other.(*mvnRef).ver) }
func (r *mvnRef) Equal(other Reference) bool       { return r.ver.Equal(other.(*mvnRef).ver) }
func (r *mvnRef) GreaterThan(other Reference) bool { return r.ver.GreaterThan(other.(*mvnRef).ver) }
func (r *mvnRef) LessThan(other Reference) bool    { return r.ver.LessThan(other.(*mvnRef).ver) }
func (r *mvnRef) GreaterThanOrEqual(other Reference) bool {
	return r.ver.GreaterThanOrEqual(other.(*mvnRef).ver)
}
func (r *mvnRef) LessThanOrEqual(other Reference) bool {
	return r.ver.LessThanOrEqual(other.(*mvnRef).ver)
}
func (r *mvnRef) Type() VersionType { return Maven }
