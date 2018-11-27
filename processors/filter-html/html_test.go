package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"bitfan/processors/doc"
	"bitfan/processors/testutils"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func getExampleConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"source_field": "source",
		"text": map[string]interface{}{
			"title":    "html head title",
			"subtitle": "html span",
			"links":    "a",
		},
		"size": map[string]interface{}{
			"comments": "html div",
			"meta":     "html meta",
		},
		"tag_on_failure": []string{"_htmlparsefailure1"},
	}
}

func TestConfigureDefaults(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{}
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.Equal(t, nil, ret, "configuration is correct, it should return nil")

	assert.Equal(t, "message", p.opt.SourceField)
	assert.Equal(t, []string{"_htmlparsefailure"}, p.opt.TagOnFailure)
}

func TestConfigure(t *testing.T) {
	p := New().(*processor)
	conf := getExampleConfiguration()
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.Equal(t, nil, ret, "configuration is correct, it should return nil")

	assert.Equal(t, "source", p.opt.SourceField)
	assert.Equal(t, []string{"_htmlparsefailure1"}, p.opt.TagOnFailure)
	assert.Equal(t, len(p.opt.Text), 3)
	assert.Equal(t, len(p.opt.Size), 2)
}
func TestReceiveInvalidHTML(t *testing.T) {
	t.Skip("TODO")
}
func TestReceiveHTML1(t *testing.T) {
	p := New()
	ctx := testutils.NewProcessorContext()
	p.Configure(ctx, getExampleConfiguration())

	em := testutils.NewPacketOld("", nil)
	em.Fields().SetValueForPath(HTML1, "source")

	p.Receive(em)

	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.Equal(t, 1, ctx.SentPacketsCount(PORT_SUCCESS), "only one event should have been sent by processor")

	em = ctx.SentPackets(PORT_SUCCESS)[0]

	assert.Equal(t, "Titre de la page", ctx.SentPackets(PORT_SUCCESS)[0].Fields().ValueOrEmptyForPathString("title"))

	v, err := em.Fields().ValueForPath("subtitle")
	assert.NoError(t, err)
	assert.Equal(t, []string{"txt1", "txt2"}, v)

	assert.Equal(t, "", em.Fields().ValueOrEmptyForPathString("links"))

	v, err = em.Fields().ValueForPath("meta")
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = em.Fields().ValueForPath("comments")
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
}

const HTML1 = `<html>
<head>
    <title>Titre de la page</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <meta description="belle description">
</head>
<body>
<span>txt1</span>
<span>txt2</span>
<!-- Ici votre site -->
</body>
</html>`
