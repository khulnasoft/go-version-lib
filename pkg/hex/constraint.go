package hex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var hexTokenRe = regexp.MustCompile(`^(>=|<=|==|!=|=|>|<|~>)?\s*(.+)$`)

var hexAndRe = regexp.MustCompile(`\band\b`)
var hexOrRe = regexp.MustCompile(`\bor\b`)

// translateHexConstraint rewrites Hex (Mix) requirement syntax into semver
// constraint syntax. Hex uses "and"/"or" keywords and the "~>" pessimistic
// operator whose upper bound depends on the number of version segments:
// "~>1.2.3" means ">=1.2.3 <1.3.0" and "~>1.2" means ">=1.2.0 <2.0.0".
func translateHexConstraint(s string) (string, error) {
	s = hexOrRe.ReplaceAllString(s, " || ")
	s = hexAndRe.ReplaceAllString(s, " ")

	var groups []string
	for _, part := range strings.Split(s, "||") {
		var toks []string
		for _, f := range strings.FieldsFunc(part, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t' || r == '\n'
		}) {
			toks = append(toks, f)
		}
		var merged []string
		for i := 0; i < len(toks); i++ {
			if isHexOp(toks[i]) && i+1 < len(toks) {
				merged = append(merged, toks[i]+toks[i+1])
				i++
				continue
			}
			merged = append(merged, toks[i])
		}
		var out []string
		for _, tok := range merged {
			t, err := translateHexToken(tok)
			if err != nil {
				return "", err
			}
			out = append(out, t)
		}
		groups = append(groups, strings.Join(out, " "))
	}
	return strings.Join(groups, " || "), nil
}

func isHexOp(s string) bool {
	switch s {
	case ">=", "<=", "==", "!=", "=", ">", "<", "~>":
		return true
	}
	return false
}

func translateHexToken(tok string) (string, error) {
	m := hexTokenRe.FindStringSubmatch(tok)
	if m == nil {
		return "", fmt.Errorf("malformed hex requirement: %q", tok)
	}
	op, ver := m[1], strings.TrimSpace(m[2])

	switch op {
	case "~>":
		return pessimisticBounds(ver)
	case "", "=", "==":
		if op == "" {
			// Hex requires an explicit operator; treat a bare version as
			// exact to stay conservative
			op = "="
		}
		return op + padHexVersion(ver), nil
	case "!=", ">", "<", ">=", "<=":
		return op + padHexVersion(ver), nil
	default:
		return "", fmt.Errorf("unsupported hex operator %q", op)
	}
}

// pessimisticBounds expands "~>V" per Hex semantics:
// 3 segments -> >=V <X.(Y+1).0 ; 2 segments -> >=V <(X+1).0.0 ;
// 1 segment -> >=V <(X+1).0.0
func pessimisticBounds(ver string) (string, error) {
	v := padHexVersion(ver)
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("malformed hex pessimistic version %q", ver)
	}
	seg := func(i int) (int, error) {
		return strconv.Atoi(parts[i])
	}
	major, err := seg(0)
	if err != nil {
		return "", fmt.Errorf("malformed hex version %q", ver)
	}
	minor, err := seg(1)
	if err != nil {
		return "", fmt.Errorf("malformed hex version %q", ver)
	}
	// count original segments to decide the bump level
	origSegs := len(strings.Split(ver, "."))
	var upper string
	switch origSegs {
	case 3:
		upper = fmt.Sprintf("%d.%d.0", major, minor+1)
	default:
		upper = fmt.Sprintf("%d.0.0", major+1)
	}
	return fmt.Sprintf(">=%s <%s", v, upper), nil
}

func padHexVersion(ver string) string {
	n := strings.Count(ver, ".")
	switch {
	case n == 0:
		return ver + ".0.0"
	case n == 1:
		return ver + ".0"
	}
	return ver
}
