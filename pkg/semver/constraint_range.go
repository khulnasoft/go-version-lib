package semver

import (
	"fmt"

	"github.com/khulnasoft/go-version-lib/pkg/part"
	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// ToVersionRange lowers the constraints to the generic versionrange model,
// enabling cross-ecosystem range algebra (Intersect, Union, Complement).
//
// Wildcard versions (e.g. "1.x", "1.2.*") are expanded to their equivalent
// lower/upper bound ranges, but only for equality-style operators; using a
// wildcard with ordering operators returns an error.
func (cs Constraints) ToVersionRange() (versionrange.Range[Version], error) {
	less := func(a, b Version) int { return a.Compare(b) }
	var sets []versionrange.Set[Version]
	for _, group := range cs.constraints {
		var cmps []versionrange.Comparator[Version]
		for _, c := range group {
			cc, err := c.toComparators()
			if err != nil {
				return versionrange.Range[Version]{}, err
			}
			cmps = append(cmps, cc...)
		}
		sets = append(sets, versionrange.NewSet(cmps...))
	}
	return versionrange.New(less, sets...), nil
}

func (c constraint) toComparators() ([]versionrange.Comparator[Version], error) {
	if wildcardIndex(c.version) >= 0 {
		switch c.op {
		case "", "=", "==":
			return WildcardComparators(c.version), nil
		default:
			return nil, fmt.Errorf("wildcard version %q only supports equality operators, got %q", c.original, c.op)
		}
	}
	cmp := func(op versionrange.Op) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{{Op: op, Value: c.version}}
	}
	between := func(upper Version) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{
			{Op: versionrange.OpGE, Value: c.version},
			{Op: versionrange.OpLT, Value: upper},
		}
	}
	switch c.op {
	case "", "=", "==":
		return cmp(versionrange.OpEQ), nil
	case "!=":
		return cmp(versionrange.OpNE), nil
	case ">":
		return cmp(versionrange.OpGT), nil
	case ">=", "=>":
		return cmp(versionrange.OpGE), nil
	case "<":
		return cmp(versionrange.OpLT), nil
	case "<=", "=<":
		return cmp(versionrange.OpLE), nil
	case "~":
		return between(c.version.TildeBump()), nil
	case "^":
		return between(c.version.CaretBump()), nil
	default:
		return nil, fmt.Errorf("unsupported operator %q", c.op)
	}
}

// wildcardIndex returns the index (0=major, 1=minor, 2=patch) of the first
// wildcard segment, or -1 if the version is fully concrete.
func wildcardIndex(v Version) int {
	segs := []part.Part{v.major, v.minor, v.patch}
	for i, p := range segs {
		// IsAny covers explicit x/*/* and default wildcard empties
		// (Empty{any:true}); zero-padded empties (Empty{any:false})
		// compare as 0 and must not be treated as wildcards.
		if p.IsAny() {
			return i
		}
	}
	return -1
}

// WildcardComparators expands a wildcard version into >=lower <upper bounds.
// A fully wildcard version ("*", "x") contributes no comparators, i.e. the
// universal set. It is exported for sibling ecosystems (e.g. npm) whose
// Version type aliases semver.Version.
func WildcardComparators(v Version) []versionrange.Comparator[Version] {
	if wildcardIndex(v) == 0 {
		return nil
	}
	lower := v.Min()
	var upper Version
	switch wildcardIndex(v) {
	case 1:
		upper = lower.IncMajor()
	case 2:
		upper = lower.IncMinor()
	default:
		// wildcard at patch or beyond with concrete major+minor
		upper = lower.IncMinor()
	}
	return []versionrange.Comparator[Version]{
		{Op: versionrange.OpGE, Value: lower},
		{Op: versionrange.OpLT, Value: upper},
	}
}
