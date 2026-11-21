package version

import "testing"

func FuzzNewVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.0.RELEASE")
	f.Add("1.0-SNAPSHOT")
	f.Add("not a version")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := NewVersion(s)
		if err != nil {
			return
		}
		// Maven's ComparableVersion semantics treat an empty version as valid.
		if s != "" && v.String() == "" {
			t.Fatalf("NewVersion(%q) produced empty String()", s)
		}
		if v.Compare(v) != 0 {
			t.Fatalf("NewVersion(%q).Compare(self) != 0", s)
		}
	})
}

func FuzzNewConstraints(f *testing.F) {
	f.Add("[1.0,2.0)")
	f.Add(">=1.0")
	f.Add("not a constraint")
	f.Fuzz(func(t *testing.T, s string) {
		if _, err := NewConstraints(s); err != nil {
			return
		}
	})
}
