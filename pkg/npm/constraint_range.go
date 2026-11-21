package npm

import (
	"fmt"

	"github.com/khulnasoft/go-version-lib/pkg/semver"
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
	if c.version.IsAny() {
		switch c.op {
		case "", "=", "==":
			return semver.WildcardComparators(c.version), nil
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
