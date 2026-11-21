package semver

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.2.3-alpha.1+build.5")
	f.Add("v0.0.0")
	f.Add("not a version")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := Parse(s)
		if err != nil {
			return
		}
		if v.String() == "" {
			t.Fatalf("Parse(%q) produced empty String()", s)
		}
		if v.Compare(v) != 0 {
			t.Fatalf("Parse(%q).Compare(self) != 0", s)
		}
	})
}

func FuzzNewConstraints(f *testing.F) {
	f.Add(">=1.2.3")
	f.Add("^1.2.3")
	f.Add("1.x || >=2.0.0 <3.0.0")
	f.Add("not a constraint")
	f.Fuzz(func(t *testing.T, s string) {
		if _, err := NewConstraints(s); err != nil {
			return
		}
	})
}
