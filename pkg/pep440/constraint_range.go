package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// ToVersionRange lowers the specifiers to the generic versionrange model,
// enabling cross-ecosystem range algebra (Intersect, Union, Complement).
//
// Specifiers with an explicit local version ("==1.0+local") and ">"
// specifiers whose version is not a post-release cannot be represented:
// PEP 440's post-release exclusion for ">" is not a contiguous interval.
// "~=" is only supported for plain release versions.
func (ss Specifiers) ToVersionRange() (versionrange.Range[Version], error) {
	// Specifiers.Check strips the candidate's local segment by default;
	// mirror that here so range satisfaction agrees with Check.
	less := func(a, b Version) int {
		return comparePublic(a, b)
	}
	var sets []versionrange.Set[Version]
	for _, group := range ss.specifiers {
		// each specifier lowers to one or more alternative conjunctions;
		// combine them pairwise (logical AND distributes over OR)
		alts := [][]versionrange.Comparator[Version]{{}}
		for _, s := range group {
			specAlts, err := s.toSets()
			if err != nil {
				return versionrange.Range[Version]{}, err
			}
			var next [][]versionrange.Comparator[Version]
			for _, a := range alts {
				for _, b := range specAlts {
					merged := append(append([]versionrange.Comparator[Version]{}, a...), b...)
					next = append(next, merged)
				}
			}
			alts = next
		}
		for _, alt := range alts {
			sets = append(sets, versionrange.NewSet(alt...))
		}
	}
	return versionrange.New(less, sets...), nil
}

// comparePublic compares two versions with local segments stripped, matching
// the default behavior of Specifiers.Check.
func comparePublic(a, b Version) int {
	ap, err := Parse(a.Public())
	if err != nil {
		return a.Compare(b)
	}
	bp, err := Parse(b.Public())
	if err != nil {
		return a.Compare(b)
	}
	return ap.Compare(bp)
}

// toSets lowers one specifier to alternative conjunctions of comparators.
func (s specifier) toSets() ([][]versionrange.Comparator[Version], error) {
	prefix := strings.HasSuffix(s.version, ".*")
	ver := strings.TrimSuffix(s.version, ".*")

	// "==1.4.*" style prefix matching is only meaningful for equality ops.
	if prefix && s.operator.op != "" && s.operator.op != "=" && s.operator.op != "==" {
		return nil, fmt.Errorf("prefix version %q only supports equality operators, got %q", s.original, s.operator.op)
	}

	cmp := func(op versionrange.Op, v Version) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{{Op: op, Value: v}}
	}
	one := func(cmps ...versionrange.Comparator[Version]) [][]versionrange.Comparator[Version] {
		return [][]versionrange.Comparator[Version]{cmps}
	}

	if prefix {
		if s.operator.op == "!=" {
			return nil, fmt.Errorf("!= with prefix version %q is not representable as a range", s.original)
		}
		lo, hi, err := prefixBounds(ver)
		if err != nil {
			return nil, err
		}
		return one(
			versionrange.Comparator[Version]{Op: versionrange.OpGE, Value: lo},
			versionrange.Comparator[Version]{Op: versionrange.OpLT, Value: hi},
		), nil
	}

	v, err := Parse(ver)
	if err != nil {
		return nil, fmt.Errorf("invalid specifier version %q: %w", s.original, err)
	}
	// Specifiers with an explicit local version ("==1.0+local") compare the
	// full local segment, while the range's less func strips locals for
	// ordinary specifiers. The two policies cannot be mixed in one range.
	if strings.Contains(s.version, "+") {
		return nil, fmt.Errorf("specifier %q with explicit local version cannot be represented as a version range", s.original)
	}

	switch s.operator.op {
	case "", "=", "==":
		return one(versionrange.Comparator[Version]{Op: versionrange.OpEQ, Value: v}), nil
	case "!=":
		return one(versionrange.Comparator[Version]{Op: versionrange.OpNE, Value: v}), nil
	case ">":
		// Check rejects post-releases of the spec's base version unless the
		// spec itself is a post-release; that exclusion is not a contiguous
		// interval and cannot be expressed as bounds.
		if !v.IsPostRelease() {
			return nil, fmt.Errorf("specifier %q cannot be represented as a version range (post-release exclusion)", s.original)
		}
		return one(versionrange.Comparator[Version]{Op: versionrange.OpGT, Value: v}), nil
	case "<":
		if v.IsPreRelease() {
			return one(versionrange.Comparator[Version]{Op: versionrange.OpLT, Value: v}), nil
		}
		// Check rejects prereleases of the spec's base version. Those occupy
		// [base.dev0, base), so <spec is (<base.dev0) OR (>=base AND <spec).
		base := v.BaseVersion()
		b, err := Parse(base)
		if err != nil {
			return nil, err
		}
		floor, err := Parse(base + ".dev0")
		if err != nil {
			return nil, fmt.Errorf("cannot construct prerelease floor for %q: %w", s.original, err)
		}
		return [][]versionrange.Comparator[Version]{
			cmp(versionrange.OpLT, floor),
			{
				{Op: versionrange.OpGE, Value: b},
				{Op: versionrange.OpLT, Value: v},
			},
		}, nil
	case ">=":
		return one(versionrange.Comparator[Version]{Op: versionrange.OpGE, Value: v}), nil
	case "<=":
		return one(versionrange.Comparator[Version]{Op: versionrange.OpLE, Value: v}), nil
	case "~=":
		return compatibleBounds(s.original, v)
	case "===":
		return nil, fmt.Errorf("arbitrary equality === is not representable as a version range: %q", s.original)
	default:
		return nil, fmt.Errorf("unsupported operator %q", s.operator.op)
	}
}

// compatibleBounds expands "~=V.N" per PEP 440: >=V.N AND ==prefix.* where
// prefix drops the last release segment. Only plain release versions are
// supported.
func compatibleBounds(original string, v Version) ([][]versionrange.Comparator[Version], error) {
	if v.IsPreRelease() || v.IsPostRelease() || v.Local() != "" || strings.Contains(original, "!") {
		return nil, fmt.Errorf("~= with %q is only supported for plain release versions", original)
	}
	release := v.BaseVersion()
	segs := strings.Split(release, ".")
	if len(segs) < 2 {
		return nil, fmt.Errorf("~= requires at least two release segments: %q", original)
	}
	prefix := strings.Join(segs[:len(segs)-1], ".")
	lo, hi, err := prefixBounds(prefix)
	if err != nil {
		return nil, err
	}
	return [][]versionrange.Comparator[Version]{
		{
			{Op: versionrange.OpGE, Value: v},
			{Op: versionrange.OpGE, Value: lo},
			{Op: versionrange.OpLT, Value: hi},
		},
	}, nil
}

// prefixBounds expands a "1.4.*"-style prefix into the interval containing
// exactly the versions whose release prefix matches, including prereleases
// of the boundary releases: [prefix.dev0, bumped.dev0).
func prefixBounds(prefix string) (Version, Version, error) {
	segs := strings.Split(prefix, ".")
	last, err := strconv.Atoi(segs[len(segs)-1])
	if err != nil {
		return Version{}, Version{}, fmt.Errorf("non-numeric prefix segment in %q", prefix)
	}
	lo, err := Parse(prefix + ".dev0")
	if err != nil {
		return Version{}, Version{}, fmt.Errorf("invalid prefix lower bound %q: %w", prefix+".dev0", err)
	}
	hiSegs := append([]string(nil), segs[:len(segs)-1]...)
	hiSegs = append(hiSegs, strconv.Itoa(last+1))
	hi, err := Parse(strings.Join(hiSegs, ".") + ".dev0")
	if err != nil {
		return Version{}, Version{}, fmt.Errorf("invalid prefix upper bound: %w", err)
	}
	return lo, hi, nil
}
