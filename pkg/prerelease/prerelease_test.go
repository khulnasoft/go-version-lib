package prerelease_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khulnasoft/go-version-lib/pkg/part"
	"github.com/khulnasoft/go-version-lib/pkg/prerelease"
)

func TestCompare(t *testing.T) {
	// both null -> equal
	assert.Equal(t, 0, prerelease.Compare(part.Parts{}, part.Parts{}))

	// any -> equal
	assert.Equal(t, 0, prerelease.Compare(part.NewParts("*"), part.NewParts("1")))
	assert.Equal(t, 0, prerelease.Compare(part.NewParts("1"), part.NewParts("*")))

	// null vs non-null: null sorts after (greater)
	assert.Equal(t, 1, prerelease.Compare(part.Parts{}, part.NewParts("alpha")))
	assert.Equal(t, -1, prerelease.Compare(part.NewParts("alpha"), part.Parts{}))

	// lexical compare of parts
	assert.Equal(t, -1, prerelease.Compare(part.NewParts("alpha"), part.NewParts("beta")))
	assert.Equal(t, 1, prerelease.Compare(part.NewParts("beta"), part.NewParts("alpha")))
	assert.Equal(t, 0, prerelease.Compare(part.NewParts("1.2"), part.NewParts("1.2")))

	// deep equal short-circuits
	assert.Equal(t, 0, prerelease.Compare(part.NewParts("1.0.0-alpha"), part.NewParts("1.0.0-alpha")))
}
