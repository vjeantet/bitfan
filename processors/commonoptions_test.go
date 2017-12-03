package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCommonOptions(t *testing.T) {
	data := getTestFields()
	AddTags([]string{"foo", "bar"}, &data)

	co := &CommonOptions{
		AddField:    map[string]interface{}{"hello": "world"},
		AddTag:      []string{"toto", "bar"},
		RemoveField: []string{"name"},
		RemoveTag:   []string{"foo"},
		Type:        "type1",
	}
	co.ProcessCommonOptions(&data)

	assert.Equal(t, "world", data.ValueOrEmptyForPathString("hello"))
	assert.False(t, data.Exists("name"))
	assert.Equal(t, "type1", data.ValueOrEmptyForPathString("type"))

	tags, _ := data.ValueForPath("tags")
	assert.Equal(t, 2, len(tags.([]string)))
	assert.Contains(t, tags, "toto")
	assert.NotContains(t, tags, "foo")
}
