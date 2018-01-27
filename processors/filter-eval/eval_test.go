package evalprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func TestMaxConcurent(t *testing.T) {
	max := New().(*processor).MaxConcurent()
	assert.Equal(t, 0, max, "this processor does support concurency")
}

func TestConfigureNoExpressionNorTemplate(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "set one expression or go template")
}

func TestReceiveSimpleMultiplication(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"expressions": map[string]interface{}{"usage": "[usage] * 100"},
		},
	)

	p.Receive(testutils.NewPacketOld("stats", map[string]interface{}{"usage": float64(1738)}))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match") {
		expected := map[string]interface{}{
			"usage": float64(173800),
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}

func TestReceiveSimpleExpressions(t *testing.T) {

	data := [][]interface{}{
		{"[usage] * 100", float64(1738), float64(173800)},
		{"[usage] + 100", float64(1738), float64(1838)},
		{"[usage] - 100", float64(1738), float64(1638)},
		{"[usage] / 100", float64(1738), float64(17.38)},
		{"[usage] & 170", 0x0F, float64(0x0A)},
		{"[usage] | 170", 0x55, float64(0xFF)},
	}

	for _, oneData := range data {
		p := New().(*processor)
		ctx := testutils.NewProcessorContext()
		p.Configure(
			ctx,
			map[string]interface{}{
				"expressions": map[string]interface{}{"usage": oneData[0]},
			},
		)

		p.Receive(testutils.NewPacketOld("stats", map[string]interface{}{"usage": oneData[1]}))
		if assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match") {
			expected := map[string]interface{}{
				"usage": oneData[2],
			}
			testutils.AssertValuesForPaths(t, ctx, expected)
		}
	}
}

func TestReceiveLeaveOtherFieldsUnchanged(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"expressions": map[string]interface{}{"usage": "[usage] * 100"},
		},
	)

	p.Receive(testutils.NewPacketOld("stats", map[string]interface{}{"usage": float64(19), "size": int64(4938), "label": "hello", "percent": float64(18.59)}))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match") {
		expected := map[string]interface{}{
			"usage":   float64(1900),
			"size":    int64(4938),
			"label":   "hello",
			"percent": float64(18.59),
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}

func TestTemplate(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"templates": map[string]interface{}{"mytest": "eval_test.tpl"},
			"var":       map[string]interface{}{"VERSION": "1.0"},
		},
	)

	fields := map[string]interface{}{"name": "Jon Doe", "templateName": "EVAL-TEST.tpl", "filter": "EVAL"}
	p.Receive(testutils.NewPacketOld("stats", fields))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match") {
		expected := map[string]interface{}{
			"mytest": "Hello Jon Doe !\n\nThis template named \"EVAL-TEST.tpl\" (version 1.0) was created for testing the \"EVAL\" filter.\n",
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}

func TestTemplateWithVar(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"templates": map[string]interface{}{"mytest": "eval_test_var.tpl"},
			"var":       map[string]interface{}{"name": "Doe Jon", "templateName": "EVAL-VAR-TEST.tpl", "filter": "EVAL", "VERSION": "1.0"},
		},
	)

	fields := map[string]interface{}{}
	p.Receive(testutils.NewPacketOld("stats", fields))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match") {
		expected := map[string]interface{}{
			"mytest": "Hello Doe Jon !\n\nThis template named \"EVAL-VAR-TEST.tpl\" (version 1.0) was created for testing the \"EVAL\" filter with var template.\n",
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}

func TestTemplateErrorWrongSyntax(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	err := p.Configure(
		ctx,
		map[string]interface{}{
			"templates": map[string]interface{}{"mytest": "eval_test_invalid.tpl"},
		},
	)

	assert.EqualError(t, err, "template: :1: function \"invalid\" not defined")

}

func TestTemplateErrorUrlNotFound(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	err := p.Configure(
		ctx,
		map[string]interface{}{
			"templates": map[string]interface{}{"mytest": "http://127.0.0.1/iuherfiuhoiuehroiuhzeroiuhaiuzheifuhaerg"},
		},
	)

	assert.EqualError(t, err, "Get http://127.0.0.1/iuherfiuhoiuehroiuhzeroiuhaiuzheifuhaerg: dial tcp 127.0.0.1:80: getsockopt: connection refused")
}

func TestTemplateErrorPathIsDirectory(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	err := p.Configure(
		ctx,
		map[string]interface{}{
			"templates": map[string]interface{}{"mytest": "/tmp"},
		},
	)

	assert.EqualError(t, err, "Error while reading \"/tmp\" [read /tmp: is a directory]")
}
