package version

import "testing"

func FuzzNewVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1:1.0-2")
	f.Add("1.0~beta1")
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
