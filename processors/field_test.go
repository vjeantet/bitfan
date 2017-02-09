package processors

import (
	"testing"

	"github.com/clbanning/mxj"
	"github.com/stretchr/testify/assert"
)

func getTestFields() mxj.Map {
	m := map[string]interface{}{
		"name": "Valere",
		"location": map[string]interface{}{
			"city":    "Paris",
			"country": "France",
		},
		"twitter": "@vjeantet",
	}
	return mxj.Map(m)
}

func TestDynamic(t *testing.T) {
	fields := getTestFields()

	str := "Hello %{name} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello Valere !", str, "")

	str = "Hello I'm %{name} I come from %{location.city} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello I'm Valere I come from Paris !", str, "")

	str = "Here nothing replaced %{unknow.path} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Here nothing replaced  !", str, "")
}
