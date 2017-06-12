package processors

import (
	"testing"
	"time"

	"github.com/clbanning/mxj"
	"github.com/stretchr/testify/assert"
)

func getTestFields() mxj.Map {

	t1, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")

	m := map[string]interface{}{
		"name": "Valere",
		"location": map[string]interface{}{
			"city":    "Paris",
			"country": "France",
		},
		"twitter":    "@vjeantet",
		"@timestamp": t1,
	}
	return mxj.Map(m)
}

func TestDynamic(t *testing.T) {
	fields := getTestFields()
	str := ""

	str = "Hello %{name} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello Valere !", str, "")

	str = "Hello I'm %{name} I come from %{location.city} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello I'm Valere I come from Paris !", str, "")

	str = "Here nothing replaced %{unknow.path} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Here nothing replaced  !", str, "")

	str = "Hello %{[name]} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello Valere !", str, "")

	str = "Hello %{[location][country]} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "Hello France !", str, "")

	str = "It's %{+YYYY.MM.dd} !"
	Dynamic(&str, &fields)
	assert.Equal(t, "It's 2012.11.01 !", str, "")

}

func BenchmarkDynamics(b *testing.B) {
	fields := getTestFields()
	str := ""
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str = "Hello %{name} !"
		Dynamic(&str, &fields)

		str = "Hello I'm %{name} I come from %{location.city} !"
		Dynamic(&str, &fields)

		str = "Here nothing replaced %{unknow} sf!"
		Dynamic(&str, &fields)

		str = "Hello %{[name]} !"
		Dynamic(&str, &fields)

		str = "Hello %{[location][country]} !"
		Dynamic(&str, &fields)

		str = "It's %{+YYYY.MM.dd} !"
		Dynamic(&str, &fields)
	}
}
