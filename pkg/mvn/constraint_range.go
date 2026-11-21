package version

import (
	"fmt"

	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// ToVersionRange lowers the constraints to the generic versionrange model,
// enabling cross-ecosystem range algebra (Intersect, Union, Complement).
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
	// Maven's ComparableVersion ordering is not a strict weak order
	// (e.g. 1.sp < 1-rc < 1 < 1.sp), so bound merging in New could drop
	// predicates that Check still enforces. Keep comparators verbatim.
	return versionrange.NewRaw(less, sets...), nil
}

func (c constraint) toComparators() ([]versionrange.Comparator[Version], error) {
	cmp := func(op versionrange.Op) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{{Op: op, Value: c.version}}
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
	default:
		return nil, fmt.Errorf("unsupported operator %q", c.op)
	}
}
