package csvcodec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"bitfan/commons"
)

func TestBuffer(t *testing.T) {
	assert.Len(t, NewDecoder(strings.NewReader("")).Buffer(), 0)
}

func TestSetOptionsError(t *testing.T) {
	d := NewDecoder(strings.NewReader("data"))
	conf := map[string]interface{}{
		"columns": 4,
	}
	err := d.SetOptions(conf, nil, "")
	assert.Error(t, err)
}

func TestDefaultSettings(t *testing.T) {
	data := "column_one,column_two,column_three\n" +
		"value   1,\"value2\",value 3\n" +
		"# Commented line\n" +
		"one, two , three,four \n"
	expectData := []map[string]interface{}{
		map[string]interface{}{
			"column_one":   "value   1",
			"column_two":   "value2",
			"column_three": "value 3",
		},
		map[string]interface{}{
			"column_one": "# Commented line",
		},
		map[string]interface{}{
			"column_one":   "one",
			"column_two":   " two ",
			"column_three": " three",
			"column4":      "four ",
		},
	}

	d := NewDecoder(strings.NewReader(data))
	var m interface{}

	for i := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err := d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestWithCustomColumns(t *testing.T) {
	data := "one two three\n" +
		"1 2\n" +
		"# Commented line\n" +
		"value1  value3 value4\n"

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"user_defined_1": "one",
			"user_defined_2": "two",
			"column3":        "three",
		},
		map[string]interface{}{
			"user_defined_1": "1",
			"user_defined_2": "2",
		},
		map[string]interface{}{
			"user_defined_1": "value1",
			"user_defined_2": "",
			"column3":        "value3",
			"column4":        "value4",
		},
	}

	d := NewDecoder(strings.NewReader(data))

	conf := map[string]interface{}{
		"autogenerate_column_names": false,
		"separator":                 " ",
		"columns":                   []string{"user_defined_1", "user_defined_2"},
		"comment":                   "#",
	}
	var l commons.Logger
	err := d.SetOptions(conf, l, "")
	assert.NoError(t, err)

	var m interface{}

	for i := range expectData {
		err = d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestMore(t *testing.T) {
	data := "column_one,column_two,column_three\n" +
		"value   1,\"value2\",value 3\n" +
		"# Commented line\n" +
		"one, two , three,four \n"
	expectData := []map[string]interface{}{
		map[string]interface{}{
			"column_one":   "value   1",
			"column_two":   "value2",
			"column_three": "value 3",
		},
		map[string]interface{}{
			"column_one": "# Commented line",
		},
		map[string]interface{}{
			"column_one":   "one",
			"column_two":   " two ",
			"column_three": " three",
			"column4":      "four ",
		},
	}

	d := NewDecoder(strings.NewReader(data))
	var m interface{}
	var i = 0
	for d.More() {
		err := d.Decode(&m)
		if i+1 <= len(expectData) {
			assert.NoError(t, err)
			assert.Equal(t, expectData[i], m)
			i = i + 1
		} else {
			assert.Error(t, err)
		}

	}
	assert.Equal(t, 3, i)
}
