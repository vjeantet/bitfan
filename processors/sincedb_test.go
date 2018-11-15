package processors

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSinceDB(t *testing.T) {
	sdboptions := &SinceDBOptions{
		WriteInterval: 1,
		Identifier:    "sincedb.json",
	}
	sdb := NewSinceDB(sdboptions)
	time.Sleep(time.Second * 2)
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
	err := sdb.Close()
	assert.Nil(t, err, "successful close of db")
}
