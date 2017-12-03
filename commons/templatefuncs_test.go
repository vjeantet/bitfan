package commons

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tf = (*templateFunctions)(nil)

func TestFuncUpper(t *testing.T) {
	for i, test := range []struct {
		s      interface{}
		expect interface{}
		isErr  bool
	}{
		{"abc", "ABC", false},
		{"123", "123", false},
		{"é", "É", false},
		{template.HTML("UpPeR"), "UPPER", false},
		{3, "3", false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := tf.upper(test.s)

		if test.isErr {
			require.Error(t, err, errMsg)
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestFuncMapValueStringOrEmpty(t *testing.T) {
	m := map[string]interface{}{
		"name": "Valere",
		"location": map[string]interface{}{
			"city":    "Paris",
			"country": "France",
		},
		"capitals": []interface {
		}{
			map[string]interface{}{"city": "Paris", "country": "France"},
			map[string]interface{}{"city": "Dakar", "country": "Sénégal"},
		},
		"twitter": "@vjeantet",
	}

	for i, test := range []struct {
		root   string
		s      string
		expect interface{}
	}{
		{"location", "city", "Paris"},
		{"location", "unknow", ""},
		{"fds", "unknow", ""},
		{"capitals", "1.city", "Dakar"},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := tf.mapValueStringOrEmpty(test.s, m[test.root])

		assert.Equal(t, test.expect, result, errMsg)
	}
}
