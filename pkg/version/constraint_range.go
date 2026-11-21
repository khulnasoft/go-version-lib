package version

import (
	"fmt"

	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// ToVersionRange lowers the constraints to the generic versionrange model,
// enabling cross-ecosystem range algebra (Intersect, Union, Complement).
//
// Note: range evaluation uses Version.Compare, which ignores build metadata
// ordering options that Constraints.Check may apply via ConstraintOption.
func (cs Constraints) ToVersionRange() (versionrange.Range[Version], error) {
	less := func(a, b Version) int { return a.Compare(b) }
	var sets []versionrange.Set[Version]
	for _, group := range cs.List() {
		var cmps []versionrange.Comparator[Version]
		for _, c := range group {
			lower, err := c.toComparators()
			if err != nil {
				return versionrange.Range[Version]{}, err
			}
			cmps = append(cmps, lower...)
		}
		sets = append(sets, versionrange.NewSet(cmps...))
	}
	return versionrange.New(less, sets...), nil
}

func (c Constraint) toComparators() ([]versionrange.Comparator[Version], error) {
	cmp := func(op versionrange.Op) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{{Op: op, Value: c.version}}
	}
	between := func(upper Version) []versionrange.Comparator[Version] {
		return []versionrange.Comparator[Version]{
			{Op: versionrange.OpGE, Value: c.version},
			{Op: versionrange.OpLT, Value: upper},
		}
	}
	switch c.operator {
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
	case "~>":
		return between(c.version.PessimisticBump()), nil
	default:
		return nil, fmt.Errorf("unsupported operator %q", c.operator)
	}
}
