package jsonlinescodec

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSettings(t *testing.T) {
	data := `{"name": "Gilbert", "wins": [["straight", "7♣"], ["one pair", "10♥"]]}
{"name": "Alexa", "wins": [["two pair", "4♠"], ["two pair", "9♠"]]}`

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"name": "Gilbert",
			"wins": []interface{}{[]interface{}{"straight", "7♣"}, []interface{}{"one pair", "10♥"}},
		},
		map[string]interface{}{
			"name": "Alexa",
			"wins": []interface{}{[]interface{}{"two pair", "4♠"}, []interface{}{"two pair", "9♠"}},
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

func TestDelimiter(t *testing.T) {
	data := `{"name": "Gilbert", "wins": [["straight", "7♣"], ["one pair", "10♥"]]}@{"name": "Alexa", "wins": [["two pair", "4♠"], ["two pair", "9♠"]]}`

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"name": "Gilbert",
			"wins": []interface{}{[]interface{}{"straight", "7♣"}, []interface{}{"one pair", "10♥"}},
		},
		map[string]interface{}{
			"name": "Alexa",
			"wins": []interface{}{[]interface{}{"two pair", "4♠"}, []interface{}{"two pair", "9♠"}},
		},
	}

	d := NewDecoder(strings.NewReader(data))
	conf := map[string]interface{}{
		"delimiter": "@",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for i := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestSetOptionsError(t *testing.T) {
	d := NewDecoder(strings.NewReader("data"))
	conf := map[string]interface{}{
		"delimiter": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}

func TestMore(t *testing.T) {
	data := `{"name": "Gilbert", "wins": [["straight", "7♣"], ["one pair", "10♥"]]}
{"name": "Alexa", "wins": [["two pair", "4♠"], ["two pair", "9♠"]]}`

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"name": "Gilbert",
			"wins": []interface{}{[]interface{}{"straight", "7♣"}, []interface{}{"one pair", "10♥"}},
		},
		map[string]interface{}{
			"name": "Alexa",
			"wins": []interface{}{[]interface{}{"two pair", "4♠"}, []interface{}{"two pair", "9♠"}},
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
		} else {
			assert.Error(t, err)
		}
		i = i + 1
	}
	assert.Equal(t, 3, i)

}
