package plaincodec

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSettings(t *testing.T) {
	data := `Stimulate carbon sunglasses garage geodesic shanty town wristwatch
skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic
 free-market apophenia. Plastic rebar semiotics wonton soup shoes   `

	d := NewDecoder(strings.NewReader(data))
	var m interface{}

	err := d.Decode(&m)
	assert.NoError(t, err)
	assert.Equal(t, data, m)

}

func TestMore(t *testing.T) {
	data := `Stimulate carbon sunglasses garage geodesic shanty town wristwatch
skyscraper. Meta-shanty town vinyl rebar claymore mine bicycle plastic
 free-market apophenia. Plastic rebar semiotics wonton soup shoes   `

	d := NewDecoder(strings.NewReader(data))

	var m interface{}
	var i = 0
	for d.More() {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, data, m)
		i = i + 1
	}
	assert.Equal(t, 1, i)
}

func TestSetOptions(t *testing.T) {
	d := NewDecoder(strings.NewReader("data"))
	conf := map[string]interface{}{
		"delimiter": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)
}
