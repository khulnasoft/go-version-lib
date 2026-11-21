package npm

import "testing"

func FuzzNewVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("v1.2.3-alpha.1")
	f.Add("not a version")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := NewVersion(s)
		if err != nil {
			return
		}
		if v.String() == "" {
			t.Fatalf("NewVersion(%q) produced empty String()", s)
		}
		if v.Compare(v) != 0 {
			t.Fatalf("NewVersion(%q).Compare(self) != 0", s)
		}
	})
}

func FuzzNewConstraints(f *testing.F) {
	f.Add(">=1.2.3")
	f.Add("^1.2.3")
	f.Add("~1.2.3")
	f.Add("not a constraint")
	f.Fuzz(func(t *testing.T, s string) {
		if _, err := NewConstraints(s); err != nil {
			return
		}
	})
}
