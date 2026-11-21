package golang_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/go-version-lib/pkg/golang"
)

func TestIsPseudoVersion(t *testing.T) {
	valid := []string{
		"v0.0.0-20191109021931-daa7c04131f5",
		"v1.2.3-pre.0.20191109021931-daa7c04131f5",
		"v1.2.4-0.20191109021931-daa7c04131f5",
	}
	for _, s := range valid {
		assert.True(t, golang.IsPseudoVersion(s), s)
	}
	invalid := []string{
		"v1.2.3",
		"1.2.3",
		"v0.0.0-20191109021931",            // missing hash
		"v0.0.0-20191109021931-XYZ123",     // bad hash
		"v0.0.0-201911090219-daa7c04131f5", // short timestamp
	}
	for _, s := range invalid {
		assert.False(t, golang.IsPseudoVersion(s), s)
	}
}

func TestParsePseudoVersion(t *testing.T) {
	pv, err := golang.ParsePseudoVersion("v1.2.4-0.20191109021931-daa7c04131f5")
	require.NoError(t, err)
	assert.Equal(t, "v1.2.4-0", pv.Base)
	assert.Equal(t, "daa7c04131f5", pv.Revision)
	assert.Equal(t, time.Date(2019, 11, 9, 2, 19, 31, 0, time.UTC), pv.Timestamp)

	_, err = golang.ParsePseudoVersion("v1.2.3")
	assert.Error(t, err)
}
