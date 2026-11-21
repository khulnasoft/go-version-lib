package versionrange_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/versionrange"
)

// intLess orders plain ints; parseInt parses comparator bounds.
func intLess(a, b int) int { return a - b }

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(s, "v"))
}

func mustParse(t *testing.T, s string) versionrange.Range[int] {
	t.Helper()
	r, err := versionrange.Parse(s, intLess, parseInt)
	require.NoError(t, err, s)
	return r
}

func TestParseAndSatisfies(t *testing.T) {
	cases := []struct {
		range_ string
		in     []int
		out    []int
	}{
		{">=1 <3", []int{0, 1, 2, 3, 4}, []int{1, 2}},
		{"1 || 3", []int{0, 1, 2, 3, 4}, []int{1, 3}},
		{"!=2 >=1 <=3", []int{0, 1, 2, 3, 4}, []int{1, 3}},
		{">= 1 , < 3", []int{0, 1, 2, 3}, []int{1, 2}},
		{"*", []int{-1, 0, 99}, []int{-1, 0, 99}},
		{">2 || <1", []int{-1, 0, 1, 2, 3}, []int{-1, 0, 3}},
		{"=5", []int{4, 5, 6}, []int{5}},
	}
	for _, tc := range cases {
		t.Run(tc.range_, func(t *testing.T) {
			r := mustParse(t, tc.range_)
			for _, v := range tc.in {
				want := false
				for _, o := range tc.out {
					if v == o {
						want = true
					}
				}
				assert.Equal(t, want, r.Satisfies(v), "range %q vs %d", tc.range_, v)
			}
		})
	}
}

func TestNormalization(t *testing.T) {
	// >=1 >=2 collapses to >=2; <=5 <=4 collapses to <=4
	r := mustParse(t, ">=1 >=2 <=5 <=4")
	assert.True(t, r.Satisfies(3))
	assert.False(t, r.Satisfies(1))
	assert.False(t, r.Satisfies(5))
	assert.True(t, r.Satisfies(4))

	// contradictory bounds -> empty range
	r = mustParse(t, ">5 <1")
	assert.True(t, r.IsEmpty())
	assert.False(t, r.Satisfies(3))

	// >1 <1 exclusive on same bound -> empty; >=1 <=1 -> single point
	r = mustParse(t, ">1 <1")
	assert.True(t, r.IsEmpty())
	r = mustParse(t, ">=1 <=1 !=1")
	assert.True(t, r.IsEmpty())
	r = mustParse(t, ">=1 <=1")
	assert.True(t, r.Satisfies(1))
	assert.False(t, r.Satisfies(2))

	// =2 absorbs other bounds
	r = mustParse(t, "=2 >=1 <=3 !=5")
	require.Len(t, r.Sets(), 1)
	require.Len(t, r.Sets()[0].Comparators(), 1)
	assert.True(t, r.Satisfies(2))
	assert.False(t, r.Satisfies(3))
}

func TestIntersect(t *testing.T) {
	a := mustParse(t, ">=1 <5")
	b := mustParse(t, ">=3")
	inter := a.Intersect(b)
	for v := 0; v < 8; v++ {
		assert.Equal(t, a.Satisfies(v) && b.Satisfies(v), inter.Satisfies(v), "v=%d", v)
	}
}

func TestUnion(t *testing.T) {
	a := mustParse(t, "<2")
	b := mustParse(t, ">5")
	u := a.Union(b)
	for v := -1; v < 8; v++ {
		assert.Equal(t, a.Satisfies(v) || b.Satisfies(v), u.Satisfies(v), "v=%d", v)
	}
}

func TestComplement(t *testing.T) {
	a := mustParse(t, ">=1 <3")
	comp := a.Complement()
	for v := -1; v < 6; v++ {
		assert.Equal(t, !a.Satisfies(v), comp.Satisfies(v), "v=%d", v)
	}

	// complement of a union
	u := mustParse(t, "<1 || >4")
	comp = u.Complement()
	for v := -1; v < 7; v++ {
		assert.Equal(t, !u.Satisfies(v), comp.Satisfies(v), "v=%d", v)
	}

	// complement of empty range is universal
	empty := mustParse(t, ">5 <1")
	comp = empty.Complement()
	assert.True(t, comp.Satisfies(0))
	assert.True(t, comp.Satisfies(100))
}

func TestString(t *testing.T) {
	r := mustParse(t, ">=1 <3")
	assert.Equal(t, ">=1 <3", r.String())
	assert.Equal(t, "(empty)", mustParse(t, ">5 <1").String())
	assert.Equal(t, "*", mustParse(t, "*").String())
}

func TestParseErrors(t *testing.T) {
	_, err := versionrange.Parse("bogus", intLess, parseInt)
	assert.Error(t, err)
	_, err = versionrange.Parse(">=abc", intLess, parseInt)
	assert.Error(t, err)
}
