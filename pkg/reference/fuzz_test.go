package reference_test

import (
	"testing"

	"github.com/khulnasoft/go-version-lib/pkg/reference"
)

func FuzzNewReference(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.0.a")
	f.Add("pkg:npm/lodash@4.17.21")
	f.Add("pkg:deb/foo/bar@1.0-2")
	f.Add("not a version")
	f.Fuzz(func(t *testing.T, s string) {
		ref, err := reference.NewReference(s)
		if err != nil {
			return
		}
		if ref.String() == "" {
			t.Fatalf("NewReference(%q) produced empty String()", s)
		}
		if ref.Compare(ref) != 0 {
			t.Fatalf("NewReference(%q).Compare(self) != 0", s)
		}
	})
}

func FuzzDetect(f *testing.F) {
	f.Add("1.2.3")
	f.Add("1.0.a")
	f.Add("pkg:npm/lodash@4.17.21")
	f.Add("not a version")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = reference.Detect(s)
	})
}
