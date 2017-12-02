package jsoncodec

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEncoderSetOptionsError(t *testing.T) {
	var b *bytes.Buffer
	d := NewEncoder(b)
	conf := map[string]interface{}{
		"Indent": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}

func TestEncoderDefaultSettings(t *testing.T) {

	data := map[string]interface{}{
		"message": "test",
		"name":    "Valere",
		"host":    "Paris",
	}
	conf := map[string]interface{}{}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)
	assert.JSONEq(t, `{"host":"Paris","message":"test","name":"Valere"}`, b.String())
}
