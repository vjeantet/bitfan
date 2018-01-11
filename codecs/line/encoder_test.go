package linecodec

import (
	"bytes"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEncoderSetOptionsError(t *testing.T) {
	var b *bytes.Buffer
	d := NewEncoder(b)
	conf := map[string]interface{}{
		"Delimiter": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}

func TestEncoderDefaultSettings(t *testing.T) {

	data := map[string]interface{}{
		"message":    "test",
		"name":       "Valere",
		"host":       "Paris",
		"@timestamp": time.Now(),
	}
	conf := map[string]interface{}{}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)
	assert.Regexp(t, `^[0-9/:]* Paris test`, b.String())
}
func TestEncoderFormat(t *testing.T) {

	data := map[string]interface{}{
		"message":    "test",
		"name":       "Valere",
		"host":       "Paris",
		"@timestamp": time.Now(),
	}
	conf := map[string]interface{}{
		"format":    "{{.host}}",
		"Delimiter": "@@",
	}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)
	assert.Regexp(t, `^Paris@@$`, b.String())
}
func TestEncoderMissingData(t *testing.T) {

	data := map[string]interface{}{
		"message": "test",
		// "city":       "Paris",
		"@timestamp": time.Now(),
	}
	conf := map[string]interface{}{
		"format": `{{.city}} - {{.message}}`,
	}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)

	assert.Equal(t, "<no value> - test\n", b.String())
}
