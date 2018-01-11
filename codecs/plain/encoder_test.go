package plaincodec

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEncoderSetOptionsError(t *testing.T) {
	var b *bytes.Buffer
	d := NewEncoder(b)
	conf := map[string]interface{}{
		"format": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}

func TestEncoderDefaultSettings(t *testing.T) {

	data := map[string]interface{}{
		"message": "test",
		"name":    "Valere",
		"city":    "Paris",
		"now":     time.Now(),
	}
	conf := map[string]interface{}{}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)

	assert.Equal(t, "test", b.String())
}

func TestEncoderFormat(t *testing.T) {
	data := map[string]interface{}{
		"message":    "test",
		"name":       "Valere",
		"city":       "Paris",
		"@timestamp": time.Now(),
	}
	conf := map[string]interface{}{
		"format": `Hello {{TS "YYYY" .}} - {{.city}} - {{.message}}`,
	}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)

	assert.Equal(t, "Hello "+strconv.Itoa(time.Now().Year())+" - Paris - test", b.String())
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

	assert.Equal(t, "<no value> - test", b.String())
}

func TestEncoderVars(t *testing.T) {

	data := map[string]interface{}{
		"message":    "test",
		"city":       "Paris",
		"@timestamp": time.Now(),
	}
	conf := map[string]interface{}{
		"var": map[string]string{
			"foo1": "bar1",
		},
		"format": `{{.city}}-${foo1:default value}-{{.message}}`,
	}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	d.SetOptions(conf, logrus.New(), "")
	d.Encode(data)

	assert.Equal(t, "Paris-bar1-test", b.String())
}

func TestEncoderFormatInvalid(t *testing.T) {
	conf := map[string]interface{}{
		"format": `{{.city}}{{if true}} - {{.message}}`,
	}
	var b = &bytes.Buffer{}
	d := NewEncoder(b)
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}
