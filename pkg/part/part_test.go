package part_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khulnasoft/go-version-lib/pkg/part"
)

func TestNewUint64(t *testing.T) {
	p, err := part.NewUint64("42")
	assert.NoError(t, err)
	assert.Equal(t, part.Uint64(42), p)
	assert.False(t, p.IsNull())
	assert.False(t, p.IsAny())
	assert.False(t, p.IsEmpty())

	_, err = part.NewUint64("abc")
	assert.Error(t, err)
}

func TestUint64Compare(t *testing.T) {
	assert.Equal(t, -1, part.Uint64(1).Compare(part.Uint64(2)))
	assert.Equal(t, 0, part.Uint64(2).Compare(part.Uint64(2)))
	assert.Equal(t, 1, part.Uint64(3).Compare(part.Uint64(2)))

	// nil compares greater
	assert.Equal(t, 1, part.Uint64(0).Compare(nil))

	// numeric vs string: string is greater
	assert.Equal(t, -1, part.Uint64(1).Compare(part.NewString("a")))
	// numeric vs pre-string: numeric is greater
	assert.Equal(t, 1, part.Uint64(1).Compare(part.NewPreString("a")))
	// numeric vs any: equal
	assert.Equal(t, 0, part.Uint64(1).Compare(part.Any(true)))
	// numeric vs empty(any): equal to zero; empty(not any): compare to zero
	assert.Equal(t, 0, part.Uint64(0).Compare(part.NewEmpty(true)))
	assert.Equal(t, 1, part.Uint64(5).Compare(part.NewEmpty(false)))
}

func TestStringPart(t *testing.T) {
	s := part.NewString("beta")
	assert.Equal(t, part.String("beta"), s)
	assert.False(t, s.IsNull())
	assert.False(t, s.IsAny())
	assert.False(t, s.IsEmpty())

	assert.Equal(t, 0, s.Compare(part.NewString("beta")))
	assert.Equal(t, -1, s.Compare(part.NewString("gamma")))
	assert.Equal(t, 1, part.NewString("gamma").Compare(s))
	// string > numeric
	assert.Equal(t, 1, s.Compare(part.Uint64(1)))
	// string == pre-string lexical compare
	assert.Equal(t, -1, s.Compare(part.NewPreString("gamma")))
}

func TestPreStringPart(t *testing.T) {
	pre := part.NewPreString("alpha")
	assert.False(t, pre.IsNull())
	assert.False(t, pre.IsAny())
	assert.False(t, pre.IsEmpty())

	// pre-string < numeric
	assert.Equal(t, -1, pre.Compare(part.Uint64(1)))
	assert.Equal(t, 0, pre.Compare(part.NewPreString("alpha")))
	assert.Equal(t, -1, pre.Compare(part.NewPreString("beta")))
}

func TestAnyPart(t *testing.T) {
	a := part.Any(true)
	assert.True(t, a.IsAny())
	assert.False(t, a.IsNull())
	assert.False(t, a.IsEmpty())
	// any compares equal to everything
	assert.Equal(t, 0, a.Compare(part.Uint64(1)))
	assert.Equal(t, 0, a.Compare(part.NewString("x")))
	assert.Equal(t, 0, a.Compare(part.NewPreString("x")))
	assert.Equal(t, 0, a.Compare(part.Any(true)))
	assert.Equal(t, 0, a.Compare(part.NewEmpty(false)))

	notAny := part.Any(false)
	assert.False(t, notAny.IsAny())
	assert.Equal(t, -1, notAny.Compare(part.Uint64(1)))

	_, err := part.NewAny("*")
	assert.NoError(t, err)
	_, err = part.NewAny("x")
	assert.NoError(t, err)
	_, err = part.NewAny("X")
	assert.NoError(t, err)
	_, err = part.NewAny("1")
	assert.Error(t, err)
}

func TestEmptyPart(t *testing.T) {
	e := part.NewEmpty(false)
	assert.True(t, e.IsEmpty())
	assert.False(t, e.IsNull())
	assert.False(t, e.IsAny())
	assert.Equal(t, 0, e.Compare(part.Uint64(0)))
	assert.Equal(t, -1, e.Compare(part.Uint64(1)))
	assert.Equal(t, 0, e.Compare(part.NewEmpty(true)))

	ea := part.NewEmpty(true)
	assert.True(t, ea.IsAny())
	assert.Equal(t, 0, ea.Compare(part.Uint64(5)))
}

func TestInfinity(t *testing.T) {
	assert.False(t, part.Infinity.IsNull())
	assert.False(t, part.Infinity.IsAny())
	assert.False(t, part.Infinity.IsEmpty())
	assert.Equal(t, 0, part.Infinity.Compare(part.Infinity))
	assert.Equal(t, 1, part.Infinity.Compare(part.Uint64(1)))
	assert.Equal(t, 1, part.Infinity.Compare(part.NewString("x")))
	assert.Equal(t, 1, part.Infinity.Compare(part.NewEmpty(true)))

	assert.False(t, part.NegativeInfinity.IsNull())
	assert.Equal(t, 0, part.NegativeInfinity.Compare(part.NegativeInfinity))
	assert.Equal(t, -1, part.NegativeInfinity.Compare(part.Uint64(1)))
	assert.Equal(t, -1, part.NegativeInfinity.Compare(part.NewString("x")))
}

func TestNewPart(t *testing.T) {
	// numeric
	p := part.NewPart("42")
	assert.Equal(t, part.Uint64(42), p)

	// wildcard
	p = part.NewPart("*")
	assert.True(t, p.(part.Any).IsAny())

	// plain string
	p = part.NewPart("beta")
	assert.Equal(t, part.String("beta"), p)
}

func TestPartsNewAndString(t *testing.T) {
	parts := part.NewParts("1.2.3")
	assert.Len(t, parts, 3)
	assert.Equal(t, "1.2.3", parts.String())

	assert.Len(t, part.NewParts(""), 0)
}

func TestPartsNormalize(t *testing.T) {
	// trailing null (zero) parts are trimmed
	parts := part.Parts{part.Uint64(1), part.Uint64(2), part.Zero}
	assert.Equal(t, part.Parts{part.Uint64(1), part.Uint64(2)}, parts.Normalize())

	// no trailing zeros
	parts = part.Parts{part.Uint64(1), part.Uint64(2), part.Uint64(3)}
	assert.Equal(t, parts, parts.Normalize())
}

func TestPartsPadding(t *testing.T) {
	parts := part.Parts{part.Uint64(1)}
	padded := parts.Padding(3, part.Zero)
	assert.Equal(t, 3, len(padded))
	assert.Equal(t, part.Uint64(1), padded[0])
	assert.Equal(t, part.Zero, padded[1])
	assert.Equal(t, part.Zero, padded[2])

	// no-op when already large enough
	assert.Equal(t, parts, parts.Padding(1, part.Zero))
}

func TestPartsCompare(t *testing.T) {
	assert.Equal(t, 0, part.NewParts("1.2.3").Compare(part.NewParts("1.2.3")))
	assert.Equal(t, -1, part.NewParts("1.2").Compare(part.NewParts("1.2.3")))
	assert.Equal(t, 1, part.NewParts("1.2.3").Compare(part.NewParts("1.2")))

	// any compares equal
	assert.Equal(t, 0, part.NewParts("1").Compare(part.NewParts("*")))
}

func TestPartsIsNullIsAny(t *testing.T) {
	assert.True(t, part.Parts{}.IsNull())
	assert.True(t, part.Parts{part.Any(true)}.IsNull())
	assert.True(t, part.Parts{part.Any(true)}.IsAny())
	assert.False(t, part.Parts{part.Uint64(1)}.IsAny())
	assert.False(t, part.Parts{}.IsAny())
	assert.False(t, part.Parts{}.IsEmpty())
}

func TestPartsZip(t *testing.T) {
	a := part.NewParts("1.2")
	b := part.NewParts("1.2.3")
	var got []struct {
		L, R interface{}
	}
	iter := a.Zip(b)
	for {
		tuple := iter()
		if tuple == nil {
			break
		}
		got = append(got, struct {
			L, R interface{}
		}{tuple.Left, tuple.Right})
	}
	assert.Len(t, got, 3)
	// third tuple: left nil, right 3
	assert.Nil(t, got[2].L)
	assert.Equal(t, part.Uint64(3), got[2].R)
}

func TestUint64SliceToParts(t *testing.T) {
	parts := part.Uint64SliceToParts([]part.Uint64{1, 2, 3})
	assert.Equal(t, part.Parts{part.Uint64(1), part.Uint64(2), part.Uint64(3)}, parts)
}
