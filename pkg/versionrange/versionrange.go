// Package versionrange provides an ecosystem-agnostic model for version
// ranges. A Range is a disjunction (logical OR) of Sets, where each Set is a
// conjunction (logical AND) of Comparators against a single version value.
//
// The zero Range matches nothing; a Set with no comparators matches every
// version (the universal set).
package versionrange

import (
	"fmt"
	"strings"
	"unicode"
)

// Op is a comparison operator applied to a Comparator's bound value.
type Op int

const (
	OpEQ Op = iota // =
	OpNE           // !=
	OpLT           // <
	OpLE           // <=
	OpGT           // >
	OpGE           // >=
)

var opStrings = map[Op]string{
	OpEQ: "=",
	OpNE: "!=",
	OpLT: "<",
	OpLE: "<=",
	OpGT: ">",
	OpGE: ">=",
}

func (o Op) String() string {
	if s, ok := opStrings[o]; ok {
		return s
	}
	return fmt.Sprintf("Op(%d)", int(o))
}

// Less compares two versions of type T, returning a negative number if a < b,
// zero if a == b, and a positive number if a > b.
type Less[T any] func(a, b T) int

// Comparator is a single bound: the constraint "Value <op> bound".
type Comparator[T any] struct {
	Op    Op
	Value T
}

// Set is a conjunction of Comparators (logical AND). Use NewSet to build
// one; the raw comparators are normalized when the set is added to a Range.
type Set[T any] struct {
	cmps  []Comparator[T]
	empty bool
}

// Comparators returns a copy of the set's comparators (nil for the
// universal set).
func (s Set[T]) Comparators() []Comparator[T] {
	return append([]Comparator[T](nil), s.cmps...)
}

// Empty reports whether the set is unsatisfiable.
func (s Set[T]) Empty() bool { return s.empty }

// NewSet returns a Set from the given comparators. The comparators are stored
// verbatim; normalization happens when the set is combined into a Range via
// New, Intersect, or Union.
func NewSet[T any](cmps ...Comparator[T]) Set[T] {
	return Set[T]{cmps: append([]Comparator[T](nil), cmps...)}
}

// Universal returns the set that matches every version.
func Universal[T any]() Set[T] { return Set[T]{} }

// Range is a disjunction of Sets (logical OR).
type Range[T any] struct {
	less Less[T]
	sets []Set[T]
	// raw disables bound merging/normalization. It is required when less is
	// not a strict weak order (e.g. Maven's ComparableVersion, which is
	// non-transitive for qualifier versions like 1.sp vs 1-rc).
	raw bool
}

// New returns a Range over the given sets, normalized: unsatisfiable sets are
// dropped and redundant comparators within each set are removed.
func New[T any](less Less[T], sets ...Set[T]) Range[T] {
	r := Range[T]{less: less}
	for _, s := range sets {
		ns := normalize(less, s)
		// drop duplicate disjuncts so Union/Intersect stay compact
		dup := false
		for _, e := range r.sets {
			if sameSet(less, e, ns) {
				dup = true
				break
			}
		}
		if !dup {
			r.sets = append(r.sets, ns)
		}
	}
	r.dropEmpty()
	return r
}

// sameSet reports whether two normalized sets contain equivalent
// comparators (same operators, pairwise-equal values under less).
func sameSet[T any](less Less[T], a, b Set[T]) bool {
	if a.empty != b.empty || len(a.cmps) != len(b.cmps) {
		return false
	}
	used := make([]bool, len(b.cmps))
	for _, ca := range a.cmps {
		found := false
		for j, cb := range b.cmps {
			if !used[j] && ca.Op == cb.Op && less(ca.Value, cb.Value) == 0 {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// NewRaw returns a Range that stores each set's comparators verbatim,
// without merging bounds or detecting contradictions. Use it when less is
// not a strict weak order, where normalization could change semantics.
// Intersect, Union, and Complement preserve the raw mode.
func NewRaw[T any](less Less[T], sets ...Set[T]) Range[T] {
	r := Range[T]{less: less, raw: true}
	for _, s := range sets {
		r.sets = append(r.sets, s)
	}
	r.dropEmpty()
	return r
}

// NewRange is an alias for New, kept for readability at call sites.
func NewRange[T any](less Less[T], sets ...Set[T]) Range[T] { return New(less, sets...) }

// Less returns the ordering function used by this range.
func (r Range[T]) Less() Less[T] { return r.less }

// Sets returns a copy of the normalized disjuncts of this range.
func (r Range[T]) Sets() []Set[T] {
	out := make([]Set[T], len(r.sets))
	for i, s := range r.sets {
		out[i] = Set[T]{cmps: append([]Comparator[T](nil), s.cmps...), empty: s.empty}
	}
	return out
}

// IsEmpty reports whether the range has no disjuncts. Note this detects
// emptiness implied by contradictory ordering bounds; it is not a proof
// that no value in the domain can satisfy the range.
func (r Range[T]) IsEmpty() bool { return len(r.sets) == 0 }

// Intersect returns the intersection of two ranges: versions must satisfy at
// least one set from each range. Both ranges must use the same ordering
// function for T; when they differ, r's ordering is used and the result may
// be incorrect.
func (r Range[T]) Intersect(o Range[T]) Range[T] {
	less := r.less
	if less == nil {
		less = o.less
	}
	res := Range[T]{less: less, raw: r.raw || o.raw}
	for _, a := range r.sets {
		for _, b := range o.sets {
			cmps := make([]Comparator[T], 0, len(a.cmps)+len(b.cmps))
			cmps = append(cmps, a.cmps...)
			cmps = append(cmps, b.cmps...)
			merged := Set[T]{cmps: cmps}
			if !res.raw {
				merged = normalize(less, merged)
			}
			// drop duplicate disjuncts so repeated intersections (e.g.
			// during Complement) don't explode combinatorially
			dup := false
			for _, e := range res.sets {
				if sameSet(less, e, merged) {
					dup = true
					break
				}
			}
			if !dup {
				res.sets = append(res.sets, merged)
			}
		}
	}
	res.dropEmpty()
	return res
}

// Union returns the union of two ranges: versions may satisfy a set from
// either range.
func (r Range[T]) Union(o Range[T]) Range[T] {
	less := r.less
	if less == nil {
		less = o.less
	}
	res := Range[T]{less: less, raw: r.raw || o.raw, sets: make([]Set[T], 0, len(r.sets)+len(o.sets))}
	res.sets = append(res.sets, r.sets...)
	res.sets = append(res.sets, o.sets...)
	if !res.raw {
		for i, s := range res.sets {
			res.sets[i] = normalize(less, s)
		}
	}
	res.dropEmpty()
	return res
}

// Satisfies reports whether v is contained in the range.
func (r Range[T]) Satisfies(v T) bool {
	for _, s := range r.sets {
		if s.empty {
			continue
		}
		ok := true
		for _, c := range s.cmps {
			if !match(r.less, c, v) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// Complement returns a range matching versions that do NOT satisfy r.
// Implemented as the union of per-set negations (each comparator inverted and
// unioned), which is exact for comparator sets over a total order.
func (r Range[T]) Complement() Range[T] {
	// !(A || B) == !A && !B; !(c1 && c2 && ...) == !c1 || !c2 || ...
	// We implement De Morgan by distributing: complement of each set is a
	// union of single-comparator negations; the complement of the range is
	// the intersection of those unions.
	var result *Range[T]
	for _, s := range r.sets {
		var neg Range[T]
		neg.less = r.less
		neg.raw = r.raw
		if s.empty {
			// an empty set contributes nothing to the union, so its
			// complement is universal — skip it in the intersection
			continue
		}
		for _, c := range s.cmps {
			neg.sets = append(neg.sets, Set[T]{cmps: []Comparator[T]{negate(c)}})
		}
		if len(s.cmps) == 0 {
			// universal set: complement is empty
			return Range[T]{less: r.less}
		}
		neg.dropEmpty()
		if result == nil {
			cp := neg
			result = &cp
		} else {
			cp := result.Intersect(neg)
			result = &cp
		}
	}
	if result == nil {
		// every disjunct was empty: the range itself is empty, so its
		// complement is the universal set
		return New(r.less, Universal[T]())
	}
	return *result
}

func (r *Range[T]) dropEmpty() {
	out := r.sets[:0]
	for _, s := range r.sets {
		if !s.empty {
			out = append(out, s)
		}
	}
	r.sets = out
}

func match[T any](less Less[T], c Comparator[T], v T) bool {
	cmp := less(v, c.Value)
	switch c.Op {
	case OpEQ:
		return cmp == 0
	case OpNE:
		return cmp != 0
	case OpLT:
		return cmp < 0
	case OpLE:
		return cmp <= 0
	case OpGT:
		return cmp > 0
	case OpGE:
		return cmp >= 0
	}
	return false
}

func negate[T any](c Comparator[T]) Comparator[T] {
	switch c.Op {
	case OpEQ:
		return Comparator[T]{Op: OpNE, Value: c.Value}
	case OpNE:
		return Comparator[T]{Op: OpEQ, Value: c.Value}
	case OpLT:
		return Comparator[T]{Op: OpGE, Value: c.Value}
	case OpLE:
		return Comparator[T]{Op: OpGT, Value: c.Value}
	case OpGT:
		return Comparator[T]{Op: OpLE, Value: c.Value}
	case OpGE:
		return Comparator[T]{Op: OpLT, Value: c.Value}
	}
	return c
}

// normalize reduces a set's comparators to a canonical minimal form:
//   - an equality comparator absorbs all others when consistent
//   - strongest lower/upper bounds are kept; dominated bounds dropped
//   - contradictory bounds mark the set empty
func normalize[T any](less Less[T], s Set[T]) Set[T] {
	if s.empty {
		return s
	}
	cmps := s.cmps

	var eqs, nes []Comparator[T]
	var lowers, uppers []Comparator[T]
	for _, c := range cmps {
		switch c.Op {
		case OpEQ:
			eqs = append(eqs, c)
		case OpNE:
			nes = append(nes, c)
		case OpGT, OpGE:
			lowers = append(lowers, c)
		case OpLT, OpLE:
			uppers = append(uppers, c)
		default:
			// unknown operator: cannot reason about it; treat the
			// conjunction as unsatisfiable rather than dropping it
			return Set[T]{empty: true}
		}
	}

	// Multiple EQs must agree, else empty.
	if len(eqs) > 0 {
		e := eqs[0]
		for _, other := range eqs[1:] {
			if less(e.Value, other.Value) != 0 {
				return Set[T]{empty: true}
			}
		}
		// EQ absorbs every bound: verify consistency, then keep just EQ
		// plus any NE that excludes a *different* value is redundant.
		for _, c := range append(append([]Comparator[T]{}, lowers...), uppers...) {
			if !match(less, c, e.Value) {
				return Set[T]{empty: true}
			}
		}
		for _, n := range nes {
			if less(n.Value, e.Value) == 0 {
				return Set[T]{empty: true}
			}
		}
		return Set[T]{cmps: []Comparator[T]{e}}
	}

	// Strongest lower bound.
	var lo *Comparator[T]
	for i := range lowers {
		c := lowers[i]
		if lo == nil || strongerLower(less, c, *lo) {
			lo = &lowers[i]
		}
	}
	// Strongest upper bound.
	var hi *Comparator[T]
	for i := range uppers {
		c := uppers[i]
		if hi == nil || strongerUpper(less, c, *hi) {
			hi = &uppers[i]
		}
	}

	// Contradiction check between best bounds.
	if lo != nil && hi != nil {
		cmp := less(lo.Value, hi.Value)
		if cmp > 0 || (cmp == 0 && (lo.Op == OpGT || hi.Op == OpLT)) {
			return Set[T]{empty: true}
		}
	}

	out := make([]Comparator[T], 0, 2+len(nes))
	if lo != nil {
		out = append(out, *lo)
	}
	if hi != nil {
		out = append(out, *hi)
	}
	out = append(out, nes...)

	// If bounds collapse to a single point, a NE on it is contradictory.
	if point, ok := singlePoint(less, out); ok {
		for _, n := range nes {
			if less(n.Value, point) == 0 {
				return Set[T]{empty: true}
			}
		}
	}
	return Set[T]{cmps: out}
}

func strongerLower[T any](less Less[T], a, b Comparator[T]) bool {
	cmp := less(a.Value, b.Value)
	if cmp != 0 {
		return cmp > 0
	}
	// same bound value: strict > beats >=
	return a.Op == OpGT && b.Op == OpGE
}

func strongerUpper[T any](less Less[T], a, b Comparator[T]) bool {
	cmp := less(a.Value, b.Value)
	if cmp != 0 {
		return cmp < 0
	}
	return a.Op == OpLT && b.Op == OpLE
}

// singlePoint reports whether the comparator list constrains the range to
// exactly one version, returning that version.
func singlePoint[T any](less Less[T], cmps []Comparator[T]) (T, bool) {
	var zero T
	for _, c := range cmps {
		if c.Op == OpEQ {
			return c.Value, true
		}
	}
	var lo, hi *Comparator[T]
	for i := range cmps {
		switch cmps[i].Op {
		case OpGT, OpGE:
			if lo == nil || strongerLower(less, cmps[i], *lo) {
				lo = &cmps[i]
			}
		case OpLT, OpLE:
			if hi == nil || strongerUpper(less, cmps[i], *hi) {
				hi = &cmps[i]
			}
		}
	}
	if lo != nil && hi != nil && lo.Op == OpGE && hi.Op == OpLE &&
		less(lo.Value, hi.Value) == 0 {
		return lo.Value, true
	}
	return zero, false
}

// String renders the range in canonical form.
func (r Range[T]) String() string {
	if r.IsEmpty() {
		return "(empty)"
	}
	parts := make([]string, 0, len(r.sets))
	for _, s := range r.sets {
		if s.empty {
			parts = append(parts, "(empty)")
			continue
		}
		if len(s.cmps) == 0 {
			parts = append(parts, "*")
			continue
		}
		cs := make([]string, 0, len(s.cmps))
		for _, c := range s.cmps {
			cs = append(cs, c.Op.String()+fmt.Sprint(c.Value))
		}
		parts = append(parts, strings.Join(cs, " "))
	}
	return strings.Join(parts, " || ")
}

// Parse parses the canonical range syntax: comparators separated by
// whitespace or commas form a conjunction; "||" separates disjuncts.
// A bare version means equality. "*", "x", or "X" alone matches everything.
//
// Examples: ">=1.2.3 <2.0.0", "1.2.3 || >=2.0.0 <3.0.0", "!=1.5.0"
func Parse[T any](s string, less Less[T], parse func(string) (T, error)) (Range[T], error) {
	if strings.TrimSpace(s) == "(empty)" {
		return New(less), nil
	}
	var sets []Set[T]
	for _, part := range strings.Split(s, "||") {
		if strings.TrimSpace(part) == "" {
			return Range[T]{}, fmt.Errorf("empty disjunct in range %q", s)
		}
		cmps, err := parseSet(part, parse)
		if err != nil {
			return Range[T]{}, err
		}
		sets = append(sets, NewSet(cmps...))
	}
	return New(less, sets...), nil
}

var cmpOps = []string{"<=", ">=", "!=", "==", "=", "<", ">"}

func parseSet[T any](s string, parse func(string) (T, error)) ([]Comparator[T], error) {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
	// Merge lone operator tokens with the following version token.
	var tokens []string
	for i := 0; i < len(fields); i++ {
		if isBareOp(fields[i]) && i+1 < len(fields) {
			tokens = append(tokens, fields[i]+fields[i+1])
			i++
			continue
		}
		tokens = append(tokens, fields[i])
	}

	var cmps []Comparator[T]
	for _, tok := range tokens {
		op, vs := splitOp(tok)
		if vs == "" {
			return nil, fmt.Errorf("empty comparator in %q", s)
		}
		if vs == "*" || vs == "x" || vs == "X" {
			if op != OpEQ {
				return nil, fmt.Errorf("wildcard %q cannot be used with operator %s", vs, op)
			}
			// wildcard alone is the universal set
			continue
		}
		v, err := parse(vs)
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", vs, err)
		}
		cmps = append(cmps, Comparator[T]{Op: op, Value: v})
	}
	return cmps, nil
}

func isBareOp(s string) bool {
	for _, op := range cmpOps {
		if s == op {
			return true
		}
	}
	return false
}

func splitOp(tok string) (Op, string) {
	for _, s := range cmpOps {
		if strings.HasPrefix(tok, s) {
			switch s {
			case "=", "==":
				return OpEQ, tok[len(s):]
			case "!=":
				return OpNE, tok[len(s):]
			case "<":
				return OpLT, tok[len(s):]
			case "<=":
				return OpLE, tok[len(s):]
			case ">":
				return OpGT, tok[len(s):]
			case ">=":
				return OpGE, tok[len(s):]
			}
		}
	}
	return OpEQ, tok
}
