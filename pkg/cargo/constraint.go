package cargo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cargoTokenRe = regexp.MustCompile(`^(>=|<=|=|\^|~|>|<|!=)?\s*(.+)$`)

// translateCargoConstraint rewrites Cargo requirement syntax into semver
// constraint syntax. Cargo treats a bare version as a caret requirement
// ("1.2.3" means ">=1.2.3, <2.0.0"), supports trailing wildcards
// ("1.2.*"), and uses commas or whitespace as AND separators.
func translateCargoConstraint(s string) (string, error) {
	var groups []string
	// split on "||" only; AND separators are commas/whitespace
	for _, part := range strings.Split(s, "||") {
		var toks []string
		for _, f := range strings.FieldsFunc(part, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t' || r == '\n'
		}) {
			toks = append(toks, f)
		}
		// merge lone operator tokens with the following version
		var merged []string
		for i := 0; i < len(toks); i++ {
			if isCargoOp(toks[i]) && i+1 < len(toks) {
				merged = append(merged, toks[i]+toks[i+1])
				i++
				continue
			}
			merged = append(merged, toks[i])
		}
		var out []string
		for _, tok := range merged {
			t, err := translateCargoToken(tok)
			if err != nil {
				return "", err
			}
			out = append(out, t)
		}
		groups = append(groups, strings.Join(out, " "))
	}
	return strings.Join(groups, " || "), nil
}

func isCargoOp(s string) bool {
	switch s {
	case ">=", "<=", "=", "^", "~", ">", "<", "!=":
		return true
	}
	return false
}

func translateCargoToken(tok string) (string, error) {
	m := cargoTokenRe.FindStringSubmatch(tok)
	if m == nil {
		return "", fmt.Errorf("malformed cargo requirement: %q", tok)
	}
	op, ver := m[1], strings.TrimSpace(m[2])

	// bare wildcard matches everything
	if ver == "*" || ver == "x" || ver == "X" {
		if op == "" {
			return "*", nil
		}
		return "", fmt.Errorf("cargo wildcard %q cannot be combined with operator %q", ver, op)
	}

	// version with trailing wildcards, e.g. "1.2.*" or "1.*"
	if strings.Contains(ver, "*") {
		if op != "" {
			return "", fmt.Errorf("cargo wildcard %q cannot be combined with operator %q", ver, op)
		}
		return wildcardBounds(ver)
	}

	switch op {
	case "":
		// bare version means caret in Cargo
		return "^" + padVersion(ver), nil
	case "=", "!=", ">", "<", ">=", "<=", "^", "~":
		return op + padVersion(ver), nil
	default:
		return "", fmt.Errorf("unsupported cargo operator %q", op)
	}
}

// wildcardBounds converts "1.2.*"-style requirements into explicit
// >=lower <upper bounds matching Cargo semantics.
func wildcardBounds(ver string) (string, error) {
	parts := strings.Split(ver, ".")
	var nums []int
	for _, p := range parts {
		if p == "*" || p == "x" || p == "X" {
			break
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return "", fmt.Errorf("malformed cargo wildcard version %q", ver)
		}
		nums = append(nums, n)
	}
	if len(nums) == 0 {
		return "*", nil
	}
	lo := make([]string, 3)
	for i := range lo {
		if i < len(nums) {
			lo[i] = strconv.Itoa(nums[i])
		} else {
			lo[i] = "0"
		}
	}
	hi := append([]int(nil), nums...)
	for len(hi) < 3 {
		hi = append(hi, 0)
	}
	// bump the last specified numeric component; for a single component
	// ("1.*") bump major, otherwise bump the component after the last
	// wildcard-free prefix's parent — i.e. treat "1.2.*" as <1.3.0
	if len(nums) == 1 {
		hi[0]++
		hi[1], hi[2] = 0, 0
	} else {
		hi[len(nums)-1]++
		for i := len(nums); i < 3; i++ {
			hi[i] = 0
		}
	}
	hiStr := make([]string, 3)
	for i, n := range hi {
		hiStr[i] = strconv.Itoa(n)
	}
	return fmt.Sprintf(">=%s <%s", strings.Join(lo, "."), strings.Join(hiStr, ".")), nil
}

// padVersion pads "1" -> "1.0.0" and "1.2" -> "1.2.0" for the semver parser.
func padVersion(ver string) string {
	if strings.ContainsAny(ver, "xX*") {
		return ver
	}
	n := strings.Count(ver, ".")
	switch {
	case n == 0:
		return ver + ".0.0"
	case n == 1:
		return ver + ".0"
	}
	return ver
}
