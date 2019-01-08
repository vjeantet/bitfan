package mutate

import (
	"testing"
	"time"

	"github.com/awillis/bitfan/processors/doc"
	"github.com/awillis/bitfan/processors/testutils"
	"github.com/clbanning/mxj"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a mutate.processos struct")
}

func getExampleConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"lowercase":    []string{"field1", "field2"},
		"uppercase":    []string{"ucfield1", "ucfield2", "ucfield3"},
		"Remove_field": []string{"rffield1", "rffield2", "rffield3", "rffield4"},
		"Add_field": map[string]interface{}{
			"adfield1": "value1",
			"adfield2": "value2",
		},
		"update": map[string]interface{}{
			"upfield1": "value3",
			"upfield2": "value4",
			"upfield3": "value5",
		},
		"Rename": map[string]interface{}{
			"rnfieldA": "rnfieldB",
		},
		"convert": map[string]interface{}{
			"fieldname": "integer",
		},
		"gsub": []string{"fngsub1", "/", "_", "fngsub2", "[\\?\\#\\-]", "."},

		"split": map[string]interface{}{
			"splitme": ",",
		},
		"strip":  []string{"trim1", "trim2"},
		"unknow": "Unknow value",

		"merge": map[string]interface{}{
			"array_dst": "array_src",
		},
	}
}

func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}

func TestConfigureError(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"update": 54,
	}
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.NotEqual(t, ret, nil, "configuration is not correct, it should return an error")
	assert.Implements(t, new(error), ret)
}

func TestConfigure(t *testing.T) {
	p := New().(*processor)
	conf := getExampleConfiguration()
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.Equal(t, ret, nil, "configuration is correct, it should return nil")

	assert.Equal(t, len(p.opt.Lowercase), 2, "lowercase options should have 2 strings")
	assert.Equal(t, len(p.opt.Uppercase), 3, "uppercase options should have 3 strings")
	assert.Equal(t, len(p.opt.Remove_field), 4, "Remove_field options should have 4 strings")
	assert.Equal(t, len(p.opt.Add_field), 2, "Add_field options should have 2 elements")
	assert.Equal(t, len(p.opt.Update), 3, "Update_field options should have 3 elements")
	assert.Equal(t, len(p.opt.Rename), 1, "Rename_field options should have 1 elements")
}

func TestReceive(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()

	p.Configure(ctx, getExampleConfiguration())

	em := testutils.NewPacketOld("test", nil)
	em.Fields().SetValueForPath("VALUE", "field1")
	em.Fields().SetValueForPath("loRem", "ucfield2")
	em.Fields().SetValueForPath("newvalue", "upfield3")
	em.Fields().SetValueForPath("myValue", "rnfieldA")
	em.Fields().SetValueForPath("4", "fieldname")
	em.Fields().SetValueForPath("abc /dEF/GHJ-K/", "fngsub1")
	em.Fields().SetValueForPath("Hello How are you ? c#omment lo-l ", "fngsub2")
	em.Fields().SetValueForPath("hello,my,name,is,yow", "splitme")

	em.Fields().SetValueForPath("bonjour\t", "trim1")
	em.Fields().SetValueForPath(" bonjour 	", "trim2")

	em.Fields().SetValueForPath([]string{"apple", "banana", "200"}, "array_dst")
	em.Fields().SetValueForPath([]string{"200", "500"}, "array_src")

	p.Receive(em)

	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.Equal(t, 1, ctx.SentPacketsCount(PORT_SUCCESS), "only one event should have been sent by processor")

	em = ctx.SentPackets(PORT_SUCCESS)[0]

	assert.Equal(t, "value1", em.Fields().ValueOrEmptyForPathString("adfield1"), "a new field should be added")
	assert.Equal(t, "value", em.Fields().ValueOrEmptyForPathString("field1"), "field's value should be lowercase")
	assert.Equal(t, "LOREM", em.Fields().ValueOrEmptyForPathString("ucfield2"), "field's value should be uppercase")
	assert.Equal(t, "value5", em.Fields().ValueOrEmptyForPathString("upfield3"), "field's value should be updated")
	assert.Equal(t, false, em.Fields().Exists("rnfieldA"), "field A should not exists")
	assert.Equal(t, true, em.Fields().Exists("rnfieldB"), "field B should exists")
	assert.Equal(t, "myValue", em.Fields().ValueOrEmptyForPathString("rnfieldB"), "field B should keep field A value")
	number, _ := em.Fields().ValueForPath("fieldname")
	assert.Equal(t, 4, number, "fieldname should be 4")

	assert.Equal(t, "abc _dEF_GHJ-K_", em.Fields().ValueOrEmptyForPathString("fngsub1"), "fngsub1 should be abc _dEF_GHJ-K_")
	assert.Equal(t, "abc _dEF_GHJ-K_", em.Fields().ValueOrEmptyForPathString("fngsub1"), "fngsub1 should be abc _dEF_GHJ-K_")
	values, _ := em.Fields().ValuesForPath("splitme")
	assert.Equal(t, []interface{}{"hello", "my", "name", "is", "yow"}, values, "split ")

	array, _ := em.Fields().ValuesForPath("array_dst")
	assert.Equal(t, []interface{}{"apple", "banana", "200", "500"}, array, "array merge")

}

func TestReceiveRemoveAllBut(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"Remove_all_but": []string{"upfield3", "field1"},
	}
	p.Configure(ctx, conf)

	em := testutils.NewPacketOld("test", nil)
	em.Fields().SetValueForPath("VALUE", "field1")
	em.Fields().SetValueForPath("loRem", "ucfield2")
	em.Fields().SetValueForPath("newvalue", "upfield3")
	em.Fields().SetValueForPath("myValue", "rnfieldA")

	p.Receive(em)

	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.Equal(t, 1, ctx.SentPacketsCount(PORT_SUCCESS), "only one event should have been sent by processor")

	em = ctx.SentPackets(PORT_SUCCESS)[0]

	assert.Equal(t, false, em.Fields().Exists("ucfield2"), "field should not exists")
	assert.Equal(t, false, em.Fields().Exists("rnfieldA"), "field should not exists")

	assert.Equal(t, true, em.Fields().Exists("field1"), "field should exists")
	assert.Equal(t, true, em.Fields().Exists("upfield3"), "field should exists")

}

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
func TestLowerCaseFields(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("TesT", "foo@vjeantet")
	data.SetValueForPath(4, "number")
	data.SetValueForPath([]interface{}{"a", "b", "c"}, "test2")
	data.SetValueForPath(map[string]interface{}{
		"o": "a1",
		"p": "b1",
	}, "map")

	options := []string{
		"name", "foo", "foo%{twitter}", "number",
	}

	LowerCaseFields(options, &data)

	assert.False(t, data.Exists("foo"))
	assert.Equal(t, "valere", data.ValueOrEmptyForPathString("name"))
	assert.Equal(t, "test", data.ValueOrEmptyForPathString("foo@vjeantet"))

	v, _ := data.ValueForPath("number")
	assert.Equal(t, 4, v)

}

func TestUpperCaseFields(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("test", "foo@vjeantet")
	data.SetValueForPath(4, "number")
	data.SetValueForPath([]interface{}{"a", "b", "c"}, "test2")
	data.SetValueForPath(map[string]interface{}{
		"o": "a1",
		"p": "b1",
	}, "map")

	options := []string{
		"name", "foo", "foo%{twitter}", "number",
	}

	UpperCaseFields(options, &data)

	assert.False(t, data.Exists("foo"))
	assert.Equal(t, "VALERE", data.ValueOrEmptyForPathString("name"))
	assert.Equal(t, "TEST", data.ValueOrEmptyForPathString("foo@vjeantet"))

	v, _ := data.ValueForPath("number")
	assert.Equal(t, 4, v)

}

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

func TestSplit(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("1 2 3 4", "message")
	data.SetValueForPath("1 2 3 4 ", "fieldB")
	data.SetValueForPath("1 2\t3 4 ", "fieldA")

	options := map[string]string{
		"message": " ",
		"fieldA":  "\t",
		"fieldB":  " ",
		"fieldC":  "\n",
	}

	Split(options, &data)

	v, _ := data.ValuesForPath("message")
	assert.Equal(t, 4, len(v))

	v, _ = data.ValuesForPath("fieldB")
	assert.Equal(t, 5, len(v))

	v, _ = data.ValuesForPath("fieldA")
	assert.Equal(t, 2, len(v))
	assert.Equal(t, "1 2", data.ValueOrEmptyForPathString("fieldA[0]"))
	assert.Equal(t, "3 4 ", data.ValueOrEmptyForPathString("fieldA[1]"))

	assert.False(t, data.Exists("fieldC"))
}

func TestJoin(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath([]string{"A", "B", "C", "D"}, "message")
	data.SetValueForPath([]interface{}{"A", "B", "C", "D"}, "fieldA")
	data.SetValueForPath([]interface{}{1, 2, 3, 4}, "fieldB")
	data.SetValueForPath([]int{1, 2, 3, 4}, "fieldC")
	data.SetValueForPath([]interface{}{1, "2", "O", 4}, "fieldD")

	options := map[string]string{
		"message": ",",
		"fieldA":  ",",
		"fieldB":  "|",
		"fieldB1": "\t",
		"fieldC":  "|",
		"fieldD":  "",
	}

	Join(options, &data)

	assert.Equal(t, "A,B,C,D", data.ValueOrEmptyForPathString("message"))
	assert.Equal(t, "A,B,C,D", data.ValueOrEmptyForPathString("fieldA"))
	assert.Equal(t, "1|2|3|4", data.ValueOrEmptyForPathString("fieldB"))
	assert.Equal(t, "1|2|3|4", data.ValueOrEmptyForPathString("fieldC"))
	assert.Equal(t, "12O4", data.ValueOrEmptyForPathString("fieldD"))

	assert.False(t, data.Exists("fieldB1"))
}

func TestConvert(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("5", "intme")
	data.SetValueForPath("5", "floatme")

	options := map[string]string{
		"intme":        "integer",
		"floatme":      "float",
		"floatmeSlice": "float",
	}

	Convert(options, &data)
	assert.Equal(t, M(5, nil), M(data.ValueForPath("intme")))
	assert.Equal(t, M(float64(5), nil), M(data.ValueForPath("floatme")))
}

func TestConvertBooleanToInteger(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath(false, "foo0")
	data.SetValueForPath(true, "foo1")
	data.SetValueForPath("0", "foo2")
	data.SetValueForPath("1", "foo3")
	data.SetValueForPath("2", "foo4")

	options := map[string]string{
		"foo0": "integer",
		"foo1": "integer",
		"foo2": "integer",
		"foo3": "integer",
		"foo4": "integer",
	}

	fixtures := []struct {
		path     string
		expected int
	}{
		{
			"foo0",
			0,
		},
		{
			"foo1",
			1,
		},
		{
			"foo2",
			0,
		},
		{
			"foo3",
			1,
		},
		{
			"foo4",
			2,
		},
	}

	Convert(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValueForPath(f.path)))
	}
}
func TestConvertStringToBoolean(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("true", "true_field")
	data.SetValueForPath("false", "false_field")
	data.SetValueForPath("True", "true_upper")
	data.SetValueForPath("False", "false_upper")
	data.SetValueForPath("1", "true_one")
	data.SetValueForPath("0", "false_zero")
	data.SetValueForPath("yes", "true_yes")
	data.SetValueForPath("no", "false_no")
	data.SetValueForPath("Y", "true_y")
	data.SetValueForPath("N", "false_n")
	data.SetValueForPath("none of the above", "wrong_field")

	options := map[string]string{
		"true_field":  "boolean",
		"false_field": "boolean",
		"true_upper":  "boolean",
		"false_upper": "boolean",
		"true_one":    "boolean",
		"false_zero":  "boolean",
		"true_yes":    "boolean",
		"false_no":    "boolean",
		"true_y":      "boolean",
		"false_n":     "boolean",
		"wrong_field": "boolean",
	}

	Convert(options, &data)
	assert.Equal(t, M(true, nil), M(data.ValueForPath("true_field")))
	assert.Equal(t, M(false, nil), M(data.ValueForPath("false_field")))
	assert.Equal(t, M(true, nil), M(data.ValueForPath("true_upper")))
	assert.Equal(t, M(false, nil), M(data.ValueForPath("false_upper")))
	assert.Equal(t, M(true, nil), M(data.ValueForPath("true_one")))
	assert.Equal(t, M(false, nil), M(data.ValueForPath("false_zero")))
	assert.Equal(t, M(true, nil), M(data.ValueForPath("true_yes")))
	assert.Equal(t, M(false, nil), M(data.ValueForPath("false_no")))
	assert.Equal(t, M(true, nil), M(data.ValueForPath("true_y")))
	assert.Equal(t, M(false, nil), M(data.ValueForPath("false_n")))
	assert.Equal(t, M("none of the above", nil), M(data.ValueForPath("wrong_field")))
}

func TestConvertNestedField(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath(map[string]interface{}{}, "foo")
	data.SetValueForPath("1000", "foo.bar")
	data.SetValueForPath(2000, "foo.bar2")

	options := map[string]string{
		"[foo][bar]": "integer",
		"foo.bar2":   "string",
	}

	Convert(options, &data)
	assert.Equal(t, M(1000, nil), M(data.ValueForPath("foo.bar")))
	assert.Equal(t, M("2000", nil), M(data.ValueForPath("foo.bar2")))
}

func TestConvertInteger(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath(123, "foo0")
	data.SetValueForPath(4203, "foo1")
	data.SetValueForPath(43, "foo2")
	data.SetValueForPath(0, "foo3")

	options := map[string]string{
		"foo0": "string",
		"foo1": "float",
		"foo2": "boolean",
		"foo3": "boolean",
	}

	fixtures := []struct {
		path     string
		expected interface{}
	}{
		{"foo0", "123"},
		{"foo1", float64(4203)},
		{"foo2", true},
		{"foo3", false},
	}

	Convert(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValueForPath(f.path)), "convert integer "+f.path)
	}
}

func TestConvertFloat(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath(float64(123.2), "foo0")
	data.SetValueForPath(float64(4203.004), "foo1")
	data.SetValueForPath(float64(43), "foo2")
	data.SetValueForPath(float64(0), "foo3")

	options := map[string]string{
		"foo0": "string",
		"foo1": "integer",
		"foo2": "boolean",
		"foo3": "boolean",
	}

	fixtures := []struct {
		path     string
		expected interface{}
	}{
		{"foo0", "123.200000"},
		{"foo1", 4203},
		{"foo2", true},
		{"foo3", false},
	}

	Convert(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValueForPath(f.path)), "convert float "+f.path)
	}
}

func TestConvertString(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("123.2", "foo0")
	data.SetValueForPath("4203", "foo1")
	data.SetValueForPath("43", "foo2")
	data.SetValueForPath("0", "foo3")
	data.SetValueForPath([]interface{}{"1", "3"}, "foo4")
	data.SetValueForPath([]string{"1", "3"}, "foo5")
	data.SetValueForPath([]interface{}{"1", 4, "lol"}, "foo6")

	options := map[string]string{
		"foo0": "float",
		"foo1": "integer",
		"foo2": "boolean",
		"foo3": "boolean",
		"foo4": "integer",
		"foo5": "integer",
		"foo6": "integer",
	}

	fixtures := []struct {
		path     string
		expected interface{}
	}{
		{"foo0", []interface{}{123.200000}},
		{"foo1", []interface{}{int(4203)}},
		{"foo2", []interface{}{"43"}},
		{"foo3", []interface{}{false}},
		{"foo4", []interface{}{1, 3}},
		{"foo5", []interface{}{1, 3}},
		{"foo6", []interface{}{1, 4, "lol"}},
	}

	Convert(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValuesForPath(f.path)), "convert string "+f.path)
	}
}

func TestConvertBoolean(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath(true, "foo0")
	data.SetValueForPath(true, "foo1")
	data.SetValueForPath(true, "foo2")

	data.SetValueForPath(false, "foo3")
	data.SetValueForPath(false, "foo4")
	data.SetValueForPath(false, "foo5")
	data.SetValueForPath([]interface{}{false, true}, "foo6")

	options := map[string]string{
		"foo0": "string",
		"foo1": "integer",
		"foo2": "float",
		"foo3": "string",
		"foo4": "integer",
		"foo5": "float",
		"foo6": "string",
	}

	fixtures := []struct {
		path     string
		expected interface{}
	}{
		{"foo0", []interface{}{"true"}},
		{"foo1", []interface{}{1}},
		{"foo2", []interface{}{float64(1)}},

		{"foo3", []interface{}{"false"}},
		{"foo4", []interface{}{0}},
		{"foo5", []interface{}{float64(0)}},

		{"foo6", []interface{}{"false", "true"}},
	}

	Convert(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValuesForPath(f.path)), "convert bool "+f.path)
	}
}

func TestMerge(t *testing.T) {
	data := getTestFields()
	data.SetValueForPath("A", "foo0")
	data.SetValueForPath([]string{"B", "C"}, "foo1")
	data.SetValueForPath([]string{"A", "B", "C"}, "foo2")
	data.SetValueForPath([]interface{}{"B", "C"}, "foo4")
	data.SetValueForPath([]interface{}{"A", "B", "C"}, "foo5")
	data.SetValueForPath([]interface{}{"D", "E"}, "foo6")
	data.SetValueForPath([]interface{}{"D", "E"}, "foo7")
	data.SetValueForPath([]interface{}{"C", "D", "E"}, "foo8")
	data.SetValueForPath([]string{"B", "C"}, "foo9")
	data.SetValueForPath([]int{1, 2, 3}, "foo10")
	data.SetValueForPath([]int{3, 4}, "foo11")
	data.SetValueForPath([]int{1, 2}, "foo12")
	data.SetValueForPath(3, "foo13")
	data.SetValueForPath([]int{1, 2}, "foo14")

	options := map[string]string{
		"foo1":  "foo0",
		"foo2":  "foo0",
		"foo4":  "foo0",
		"foo5":  "foo0",
		"foo6":  "foo5",
		"foo7":  "foo9",
		"foo8":  "foo2",
		"foo10": "foo9",
		"foo12": "foo11",
		"foo14": "foo13",
		"fooZ":  "foo13",
		"foo15": "fooY",
	}

	fixtures := []struct {
		path     string
		expected interface{}
	}{
		{"foo0", []interface{}{"A"}},
		{"foo1", []interface{}{"B", "C", "A"}},
		{"foo2", []interface{}{"A", "B", "C"}},
		{"foo4", []interface{}{"B", "C", "A"}},
		{"foo5", []interface{}{"A", "B", "C"}},
		{"foo6", []interface{}{"D", "E", "A", "B", "C"}},
		{"foo7", []interface{}{"D", "E", "B", "C"}},
		{"foo8", []interface{}{"C", "D", "E", "A", "B"}},
		{"foo10", []interface{}{1, 2, 3, "B", "C"}},
		{"foo12", []interface{}{1, 2, 3, 4}},
		{"foo14", []interface{}{1, 2, 3}},
	}

	Merge(options, &data)

	for _, f := range fixtures {
		assert.Equal(t, M(f.expected, nil), M(data.ValuesForPath(f.path)), "convert merge "+f.path)
	}
	assert.False(t, false, data.Exists("foo15"))
	assert.False(t, false, data.Exists("fooZ"))
}

// https://bitfan/issues/71
func TestNoEmptyTags(t *testing.T) {
	Convey("When no option about tags is involved", t, func() {
		event := testutils.NewPacketOld("", map[string]interface{}{})
		conf := map[string]interface{}{
			"target": "name1",
		}
		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)

		Convey("Then no empty tags fields is created in resulting event", func() {
			em := p.SentPackets(0)[0]
			So(em.Fields().Exists("tags"), ShouldBeFalse)

		})
	})

}

//Shim for 2 param return values
func M(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}
