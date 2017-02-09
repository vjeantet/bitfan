package mutate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
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

	em := testutils.NewPacket("test", nil)
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
	value, _ := em.Fields().ValueForPath("splitme")
	assert.Equal(t, []string{"hello", "my", "name", "is", "yow"}, value, "split ")

	array, _ := em.Fields().ValueForPath("array_dst")
	assert.Equal(t, []string{"apple", "banana", "200", "500"}, array, "array merge")

}

func TestReceiveRemoveAllBut(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"Remove_all_but": []string{"upfield3", "field1"},
	}
	p.Configure(ctx, conf)

	em := testutils.NewPacket("test", nil)
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
