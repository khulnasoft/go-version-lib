package version

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.0a1")
	f.Add("1.0.post1.dev2")
	f.Add("v1!2.0")
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

func FuzzNewSpecifiers(f *testing.F) {
	f.Add(">=1.2.3")
	f.Add("~=1.4.5")
	f.Add(">=1.0,<2.0")
	f.Add("not a specifier")
	f.Fuzz(func(t *testing.T, s string) {
		if _, err := NewSpecifiers(s); err != nil {
			return
		}
	})
}
