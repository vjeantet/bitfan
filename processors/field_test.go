package processors

import (
	"testing"
	"time"

	"github.com/clbanning/mxj"
	"github.com/k0kubun/pp"
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

func TestFieldSetType(t *testing.T) {
	data := getTestFields()
	SetType("Type2", &data)
	assert.Equal(t, "Type2", data.ValueOrEmptyForPathString("type"))
}
func TestFieldSetTypeDynamic(t *testing.T) {
	data := getTestFields()
	SetType("V%{[name]}", &data)
	assert.Equal(t, "VValere", data.ValueOrEmptyForPathString("type"))
}
func TestFieldSetTypeEmptyDontOverWriteExistingOne(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("Type1", "type")
	SetType("", &data)
	assert.Equal(t, "Type1", data.ValueOrEmptyForPathString("type"))
}

func TestFieldAddFields(t *testing.T) {
	data := getTestFields()
	add_fields := map[string]interface{}{
		"foo1":          "bar1",
		"foo2":          "bar2",
		"foo3Dyn":       "%{name}-ok",
		"foo%{twitter}": "twitOK",
		"test": map[string]interface{}{
			"o": "v",
		},
		"test2": []interface{}{"A", "B"},
	}
	AddFields(add_fields, &data)
	assert.Equal(t, "bar1", data.ValueOrEmptyForPathString("foo1"))
	assert.Equal(t, "bar2", data.ValueOrEmptyForPathString("foo2"))
	assert.Equal(t, "Valere-ok", data.ValueOrEmptyForPathString("foo3Dyn"))
	assert.Equal(t, "twitOK", data.ValueOrEmptyForPathString("foo@vjeantet"))
	assert.Equal(t, "A", data.ValueOrEmptyForPathString("test2[0]"))
	assert.Equal(t, "v", data.ValueOrEmptyForPathString("test.o"))
}
func TestFieldAddFieldsDontOverwriteExistingOnes(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("1rab", "foo1")
	add_fields := map[string]interface{}{
		"foo1": "bar1",
		"foo2": "bar2",
	}
	AddFields(add_fields, &data)
	assert.Equal(t, "1rab", data.ValueOrEmptyForPathString("foo1"))
	assert.Equal(t, "bar2", data.ValueOrEmptyForPathString("foo2"))
}

func TestFieldAddTags(t *testing.T) {
	data := getTestFields()
	add_tags := []string{"foo", "bar", "foo%{twitter}"}

	AddTags(add_tags, &data)

	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.Contains(t, tags, "foo")
	assert.Contains(t, tags, "bar")
	assert.Contains(t, tags, "foo@vjeantet")
}
func TestFieldAddTagsWithExistingOneAsString(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("toto", "tags")
	add_tags := []string{"foo", "bar", "foo%{twitter}"}

	AddTags(add_tags, &data)
	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.Contains(t, tags, "toto")
	assert.Contains(t, tags, "foo")
}
func TestFieldAddTagsWithExistingOneAsInterfaceArray(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath([]interface{}{"toto"}, "tags")
	add_tags := []string{"foo", "bar", "foo%{twitter}"}

	AddTags(add_tags, &data)
	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.Contains(t, tags, "toto")
	assert.Contains(t, tags, "foo")
}
func TestFieldAddTagsWithDontDuplicate(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath([]string{"bar"}, "tags")
	add_tags := []string{"foo", "bar", "foo%{twitter}"}

	AddTags(add_tags, &data)
	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.Contains(t, tags, "bar")
	assert.Contains(t, tags, "foo")
	assert.Equal(t, 3, len(tags.([]string)))
}
func TestFieldAddTagsWithExistingOnes(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath([]string{"toto"}, "tags")
	add_tags := []string{"foo", "bar", "foo%{twitter}"}

	AddTags(add_tags, &data)

	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.Contains(t, tags, "toto")
	assert.Contains(t, tags, "foo")
}

func TestFieldRemoveTags(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath([]string{"foo", "bar", "toto"}, "tags")

	remove_tags := []string{"bar", "foo%{twitter}"}

	RemoveTags(remove_tags, &data)

	tags, err := data.ValueForPath("tags")
	assert.NoError(t, err)
	assert.NotContains(t, tags, "bar")
	assert.Equal(t, 2, len(tags.([]string)))
}
func TestFieldRemoveTagsWhenFieldDoesNotExist(t *testing.T) {
	data := getTestFields()

	remove_tags := []string{"bar", "foo%{twitter}"}

	RemoveTags(remove_tags, &data)
	assert.False(t, data.Exists("tags"))
}

func TestFieldRemoveFields(t *testing.T) {
	data := getTestFields()

	remove_fields := []string{"name", "unknow", "twitter"}

	RemoveFields(remove_fields, &data)

	assert.False(t, data.Exists("name"))
	assert.False(t, data.Exists("twitter"))
}

func TestFieldRemoveAllButFields(t *testing.T) {
	data := getTestFields()

	but_fields := []string{"name", "unknow", "twitter"}

	RemoveAllButFields(but_fields, &data)

	assert.False(t, data.Exists("location"))
	assert.False(t, data.Exists("@timestamp"))
	assert.Equal(t, "Valere", data.ValueOrEmptyForPathString("name"))
	assert.Equal(t, "@vjeantet", data.ValueOrEmptyForPathString("twitter"))
}

// m := map[string]interface{}{
// 	"name": "Valere",
// 	"location": map[string]interface{}{
// 		"city":    "Paris",
// 		"country": "France",
// 	},
// 	"twitter":    "@vjeantet",
// 	"@timestamp": t1,
// }
func TestFieldRenameFields(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("test", "foo@vjeantet")
	data.SetValueForPath("test", "foo3Dyn")
	data.SetValueForPath([]interface{}{"A", "B", "C"}, "test2")
	data.SetValueForPath(map[string]interface{}{
		"o": "A1",
		"p": "B1",
	}, "map")

	options := map[string]string{
		"name":          "nom",
		"foo":           "bar",
		"test2":         "test3",
		"map":           "plan",
		"foo3Dyn":       "%{location.city}-ok",
		"foo%{twitter}": "twitOK",
	}

	RenameFields(options, &data)
	pp.Println("data-->", data)

	assert.False(t, data.Exists("name"))
	assert.True(t, data.Exists("nom"))
	assert.Equal(t, "Valere", data.ValueOrEmptyForPathString("nom"))

	assert.False(t, data.Exists("foo"))

	assert.False(t, data.Exists("test2"))
	assert.True(t, data.Exists("test3"))

	assert.False(t, data.Exists("map"))
	assert.True(t, data.Exists("plan"))
	v, _ := data.ValueForPath("plan")
	assert.Equal(t, 2, len(v.(map[string]interface{})))

	assert.False(t, data.Exists("foo3Dyn"))
	assert.True(t, data.Exists("Paris-ok"))

	assert.False(t, data.Exists("foo@vjeantet"))
	assert.True(t, data.Exists("twitOK"))
}

func TestFieldUpdateFields(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("test", "foo@vjeantet")
	data.SetValueForPath("test", "foo3Dyn")
	data.SetValueForPath([]interface{}{"A", "B", "C"}, "test2")
	data.SetValueForPath(map[string]interface{}{
		"o": "A1",
		"p": "B1",
	}, "map")

	options := map[string]interface{}{
		"name":          "Alex",
		"foo":           "bar",
		"foo3Dyn":       "%{location.city}-ok",
		"foo%{twitter}": "twitOK",
		"map": map[string]interface{}{
			"o": "B2",
		},
		"test2": []interface{}{"D", "E"},
	}

	UpdateFields(options, &data)

	assert.False(t, data.Exists("foo"))
	assert.Equal(t, "Alex", data.ValueOrEmptyForPathString("name"))
	assert.Equal(t, "twitOK", data.ValueOrEmptyForPathString("foo@vjeantet"))
	assert.Equal(t, "Paris-ok", data.ValueOrEmptyForPathString("foo3Dyn"))
	assert.Equal(t, "Paris", data.ValueOrEmptyForPathString("location.city"))

	v, _ := data.ValuesForPath("test2")
	assert.Equal(t, []interface{}{"D", "E"}, v)

	assert.False(t, data.Exists("map.p"))
	assert.Equal(t, "B2", data.ValueOrEmptyForPathString("map.o"))
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
