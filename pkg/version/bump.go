package version

// BumpKind identifies which version component to increment.
type BumpKind int

const (
	BumpMajor BumpKind = iota
	BumpMinor
	BumpPatch
)

// Bump returns a new Version with the given component incremented and all
// lower-order components reset to zero. Pre-release and build metadata are
// cleared. Versions with fewer segments than the bump target are treated as
// if padded with zeros (e.g. bumping patch on "1.2" yields "1.2.1").
func (v Version) Bump(kind BumpKind) Version {
	v = v.copy()
	for len(v.segments) < 3 {
		v.segments = append(v.segments, 0)
	}
	switch kind {
	case BumpMajor:
		v.segments[0]++
		v.segments[1], v.segments[2] = 0, 0
	case BumpMinor:
		v.segments[1]++
		v.segments[2] = 0
	case BumpPatch:
		v.segments[2]++
	}
	if len(v.segments) > 3 {
		for i := 3; i < len(v.segments); i++ {
			v.segments[i] = 0
		}
	}
	v.preRelease = nil
	v.buildMetadata = ""
	v.original = v.String()
	return v
}
