package rubydebugcodec

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEncoderDefaultSettings(t *testing.T) {

	data := map[string]interface{}{
		"message": "test",
		"name":    "Valere",
		"host":    "Paris",
	}
	conf := map[string]interface{}{}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)
	err = d.Encode(data)
	assert.NoError(t, err)
}
