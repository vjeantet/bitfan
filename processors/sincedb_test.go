package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSinceDB(t *testing.T) {
	sdboptions := &SinceDBOptions{
		Identifier: "sincedb.json",
	}
	sdb := NewSinceDB(sdboptions)
	assert.IsType(t, (*SinceDB)(nil), sdb)
	assert.False(t, sdb.dryrun)
}

func TestNewSinceDBDryRun(t *testing.T) {
	sdboptions := &SinceDBOptions{
		Identifier: "/dev/null",
	}
	sdb := NewSinceDB(sdboptions)
	assert.IsType(t, (*SinceDB)(nil), sdb)
	assert.True(t, sdb.dryrun)
}
