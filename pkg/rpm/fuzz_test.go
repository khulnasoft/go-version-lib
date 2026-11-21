package version

import "testing"

func FuzzNewVersion(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1:1.0-2")
	f.Add("1.0^1")
	f.Add("0:1.0~rc1-1")
	f.Fuzz(func(t *testing.T, s string) {
		v := NewVersion(s)
		// rpm.NewVersion is intentionally permissive and never fails; degenerate
		// inputs (e.g. "", ":") legitimately yield an empty version string.
		if v.Compare(v) != 0 {
			t.Fatalf("NewVersion(%q).Compare(self) != 0", s)
		}
	})
}
