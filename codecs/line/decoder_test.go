package linecodec

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assert.Len(t, NewDecoder(strings.NewReader("")).Buffer(), 0)
}

func TestDefaultSettings(t *testing.T) {
	data := `Stimulate carbon sunglasses garage geodesic shanty town wristwatch
skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic
 free-market apophenia. Plastic rebar semiotics wonton soup shoes   `

	expectData := []string{
		"Stimulate carbon sunglasses garage geodesic shanty town wristwatch",
		"skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic",
		" free-market apophenia. Plastic rebar semiotics wonton soup shoes   ",
	}

	d := NewDecoder(strings.NewReader(data))
	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err := d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestDelimiter(t *testing.T) {
	data := `Stimulate carbon sunglasses garage
	 geodesic shanty town wristwatch@skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic@ free-market apophenia. Plastic rebar semiotics wonton soup shoes   `

	expectData := []string{
		"Stimulate carbon sunglasses garage\n\t geodesic shanty town wristwatch",
		"skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic",
		" free-market apophenia. Plastic rebar semiotics wonton soup shoes   ",
	}

	d := NewDecoder(strings.NewReader(data))
	conf := map[string]interface{}{
		"delimiter": "@",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
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
	data := `Stimulate carbon sunglasses garage geodesic shanty town wristwatch
skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic
 free-market apophenia. Plastic rebar semiotics wonton soup shoes   `

	expectData := []string{
		"Stimulate carbon sunglasses garage geodesic shanty town wristwatch",
		"skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic",
		" free-market apophenia. Plastic rebar semiotics wonton soup shoes   ",
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
