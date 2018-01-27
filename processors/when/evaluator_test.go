package when

import (
	"testing"

	"golang.org/x/sync/syncmap"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
)

func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func TestExpressionEvaluation(t *testing.T) {
	event := newTestEvent()

	checkTrue(t, event, `[testInt] == 4`)
	checkTrue(t, event, `[testInt] == 8/2`)
	checkTrue(t, event, `[testInt] == [testInt3]+1`)
	checkTrue(t, event, `!(false)`)
	checkTrue(t, event, `"_grokparsefailure" in [tags]`)
	checkTrue(t, event, `! ( '_mumu' in [tags] )`)
	checkTrue(t, event, `!("_mumu" in [tags]) && [way] == "SEND"`)
	checkTrue(t, event, `[testString] == "true"`)
	checkTrue(t, event, `[location.city] == "Paris"`)
	checkTrue(t, event, `[testInt] == 4`)
	checkTrue(t, event, `[way]`)
	checkTrue(t, event, `[testInt] > [testInt3]`)
	checkTrue(t, event, `true && true`)
	checkTrue(t, event, `true || false`)
	checkTrue(t, event, `"foo" in ("foor", "foo","bar")`)
	checkTrue(t, event, `!("foo"  in ("foos", "sfoo","bar"))`)
	checkTrue(t, event, `[way] =~ '(RECEIVE|SEND)'`)
	checkFalse(t, event, `!(true)`)
	checkFalse(t, event, `"grokparsefailure" in [tags]`)
	checkFalse(t, event, `"_mumu" in [tags]`)
	checkFalse(t, event, `!("_grokparsefailure"  in [tags]) || [way] != "SEND"`)
	checkFalse(t, event, `[testString] == "false"`)
	checkFalse(t, event, `[testString] != "true"`)
	checkFalse(t, event, `[location.city] != "Paris"`)
	checkFalse(t, event, `[testInt] > 30`)
	checkFalse(t, event, `true && false`)
	checkFalse(t, event, `false || false`)
	checkFalse(t, event, `"foo" in ("foor", "foos","bar")`)
	checkFalse(t, event, `!("foo"  in ("foo", "sfoo","bar"))`)
	checkFalse(t, event, `[testUnk] == 3`)
	checkFalse(t, event, `[testUnk]`)
	checkFalse(t, event, `[way] !~ '(RECEIVE|SEND)'`)

	checkError(t, event, `[testUnk > 3`)
	checkError(t, event, `[testUnk] > 3`)
	checkError(t, event, ``)
}

func newTestEvent() processors.IPacket {
	m := map[string]interface{}{
		"testString":  "true",
		"testYes":     "yes",
		"testY":       "y",
		"testNo":      "no",
		"testN":       "n",
		"test1String": "1",
		"test1Int":    1,
		"test0String": "0",
		"test0Int":    0,
		"testBool":    true,
		"testInt":     4,
		"testInt3":    3,
		"way":         "SEND",
		"name":        "Valere",
		"tags": []string{
			"mytag",
			"_grokparsefailure",
			"_dateparsefailure",
		},
		"location": map[string]interface{}{
			"city":    "Paris",
			"country": "France",
		},
	}

	return testutils.NewPacketOld("test", m)
}

func checkError(t *testing.T, event processors.IPacket, expression string) {
	p := &processor{compiledExpressions: &syncmap.Map{}}
	_, err := p.assertExpressionWithFields(0, expression, event)
	assert.Error(t, err, expression)
}

func checkTrue(t *testing.T, event processors.IPacket, expression string) {
	p := &processor{compiledExpressions: &syncmap.Map{}}
	result, err := p.assertExpressionWithFields(0, expression, event)
	assert.NoError(t, err, "err is not nil")
	assert.True(t, result, expression)
}

func checkFalse(t *testing.T, event processors.IPacket, expression string) {
	p := &processor{compiledExpressions: &syncmap.Map{}}
	result, err := p.assertExpressionWithFields(0, expression, event)
	assert.NoError(t, err, "err is not nil")
	assert.False(t, result, expression)
}
