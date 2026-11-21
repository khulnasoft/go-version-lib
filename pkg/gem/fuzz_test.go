package gem

import "testing"

func FuzzNewVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.0.a")
	f.Add("1.2.3.pre.1")
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
	f.Add(">= 1.2.3")
	f.Add("~> 1.2")
	f.Add("not a constraint")
	f.Fuzz(func(t *testing.T, s string) {
		if _, err := NewConstraints(s); err != nil {
			return
		}
	})
}
