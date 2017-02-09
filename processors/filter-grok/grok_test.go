package grok

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
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func getExampleConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"remove_field": []string{"rffield1", "rffield2", "rffield3", "rffield4"},
		"add_field": map[string]interface{}{
			"adfield1": "value1",
			"adfield2": "value2",
		},
		"match": map[string]interface{}{
			"message":   "%{COMMONAPACHELOG}",
			"timestamp": "%{MONTHDAY:jour}/%{MONTH:mois}/%{YEAR:annee}",
		},
		"named_captures_only": true,
		"unknow":              "Unknow value",
	}
}

func TestConfigureError(t *testing.T) {
	p := New().(*processor)

	conf := map[string]interface{}{
		"match": 54,
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
	assert.Equal(t, ret, nil, "configuration is correct, error should be nil")

	assert.Equal(t, len(p.opt.AddField), 2, "Add_field options should have 2 elements")
	assert.Equal(t, len(p.opt.RemoveField), 4, "Remove_field options should have 4 elements")
	assert.Equal(t, len(p.opt.Match), 2, "Match options should have 2 elements")
	assert.Equal(t, p.opt.NamedCapturesOnly, true, "Named_captures_only should be true")
}

func TestReceive(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match": map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)

	p.Receive(em)
	assert.Equal(t, "evita", em.Fields().ValueOrEmptyForPathString("logsource"), "field value not proprely groked")
	assert.Equal(t, "connect from camomile.cloud9.net[168.100.1.3]", em.Fields().ValueOrEmptyForPathString("message"), "field value not proprely groked")
}

func TestReceiveFailure(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match": map[string]interface{}{"message": `%{UNKNOW}`},
	}

	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("hello world", nil)
	em.Fields().SetValueForPath("VALUE", "field1")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)
	// pp.Print(em.Fields())
	assert.Equal(t, "VALUE", em.Fields().ValueOrEmptyForPathString("field1"), "field value should stay")

	tags, _ := em.Fields().ValueForPath("tags")
	assert.Contains(t, tags.([]string), "_grokparsefailure", "failure tag should be set")

}

func TestRemoveTagNoTags(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":      map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"remove_tag": []string{"field1"},
	}

	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	// em.Fields().SetValueForPath([]string{"myTag", "field1", "myTag2"}, "notags")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	_, err := em.Fields().ValueForPath("tags")
	assert.NotNil(t, err, "...")
}

func TestRemoveTag(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":      map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"remove_tag": []string{"field1"},
	}

	p.Configure(ctx, conf)
	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	em.Fields().SetValueForPath([]string{"myTag", "field1", "myTag2"}, "tags")
	// em.Fields().SetValueForPath("newvalue", "upfield3")
	// em.Fields().SetValueForPath("myValue", "rnfieldA")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	assert.Equal(t, "evita", em.Fields().ValueOrEmptyForPathString("logsource"), "field value not proprely groked")

	tags, _ := em.Fields().ValueForPath("tags")
	assert.Len(t, tags.([]string), 2, "2 tags should be found")
	assert.Contains(t, tags.([]string), "myTag", "failure tag should be set")
	assert.Contains(t, tags.([]string), "myTag2", "failure tag should be set")
	assert.NotContains(t, tags.([]string), "field1", "failure tag should be set")
}

func TestAddTagToNoTags(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":   map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"add_tag": []string{"tag1", "tag2"},
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	tags, err := em.Fields().ValueForPath("tags")
	assert.Nil(t, err, "...")
	assert.Len(t, tags.([]string), 2, "2 tags should be found")
}

func TestAddTag(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":   map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"add_tag": []string{"tiptop", "tiptop2"},
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	em.Fields().SetValueForPath([]string{"myTag", "field1", "myTag2"}, "tags")
	// em.Fields().SetValueForPath("newvalue", "upfield3")
	// em.Fields().SetValueForPath("myValue", "rnfieldA")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	assert.Equal(t, "evita", em.Fields().ValueOrEmptyForPathString("logsource"), "field value not proprely groked")

	tags, _ := em.Fields().ValueForPath("tags")

	assert.Len(t, tags.([]string), 5, "5 tags should be found")

}

func TestRemoveField(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":        map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"remove_field": []string{"field1"},
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	em.Fields().SetValueForPath("valueA", "field1")
	em.Fields().SetValueForPath("valueB", "field2")
	// em.Fields().SetValueForPath("newvalue", "upfield3")
	// em.Fields().SetValueForPath("myValue", "rnfieldA")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	assert.Equal(t, "evita", em.Fields().ValueOrEmptyForPathString("logsource"), "field value not proprely groked")

	assert.False(t, em.Fields().Exists("field1"), "field1 should be removed")
	assert.True(t, em.Fields().Exists("field2"), "field2 should exists")
	assert.Equal(t, "valueB", em.Fields().ValueOrEmptyForPathString("field2"), "field2's should remain unchanged")
}
func TestAddField(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":     map[string]interface{}{"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`},
		"add_field": map[string]interface{}{"field1": `Hello World`},
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	em.Fields().SetValueForPath("valueB", "field2")
	// em.Fields().SetValueForPath("newvalue", "upfield3")
	// em.Fields().SetValueForPath("myValue", "rnfieldA")

	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	assert.Equal(t, "evita", em.Fields().ValueOrEmptyForPathString("logsource"), "field value not proprely groked")

	assert.True(t, em.Fields().Exists("field2"), "field2 should exists")
	assert.Equal(t, "Hello World", em.Fields().ValueOrEmptyForPathString("field1"), "field's should remain unchanged")

}

func TestPatterns_dirError(t *testing.T) {
	p := New().(*processor)

	conf := getExampleConfiguration()
	conf["patterns_dir"] = []string{"/tmp/unknow"}
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)

	assert.NotEqual(t, ret, nil, "configuration is not correct, error should not be nil")

}

func TestTag_on_failure(t *testing.T) { t.Skip("...") }

func TestNamed_captures_only(t *testing.T) { t.Skip("...") }

func TestKeep_empty_captures(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.opt.Match = map[string]interface{}{
		"message": `%{COMBINEDAPACHELOG}`,
	}
	p.Configure(ctx, map[string]interface{}{"keep_empty_captures": true})

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket(`127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"`, nil)
	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	assert.Equal(t, "", em.Fields().ValueOrEmptyForPathString("rawrequest"), "field value not proprely groked")
}

func TestKeep_empty_capturesFalse(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.opt.Match = map[string]interface{}{
		"message": `%{COMBINEDAPACHELOG}`,
	}
	p.Configure(ctx, map[string]interface{}{"keep_empty_captures": false})

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket(`127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"`, nil)
	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	_, err := em.Fields().ValueForPath("rawrequest")
	assert.NotNil(t, err, "field should not exists")
}

func TestBreak_on_matchFalse(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match": map[string]interface{}{
			"unknow":  `%{NUMBER} %{GREEDYDATA:message}`,
			"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`,
			"program": `%{WORD:programname}/%{WORD:daemon}`,
		},
		"break_on_match": false,
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)
	// pp.Print(em.Fields())
	assert.Equal(t, "smtpd", em.Fields().ValueOrEmptyForPathString("daemon"), "field value not proprely groked")
}

func TestBreak_on_matchTrue(t *testing.T) {
	p := New().(*processor)

	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match": map[string]interface{}{
			"unknow":  `%{NUMBER} %{GREEDYDATA:message}`,
			"message": `%{SYSLOGBASE} %{GREEDYDATA:message}`,
			"program": `%{GREEDYDATA:programname}/%{GREEDYDATA:daemon}`,
		},
		"break_on_match": true,
	}
	p.Configure(ctx, conf)

	//NewTestEvent(sourceAgentName string, message string, fields map[string]interface{}) Event {
	em := testutils.NewPacket("Mar 16 00:01:25 evita postfix/smtpd[1713]: connect from camomile.cloud9.net[168.100.1.3]", nil)
	// em.On("Pipe", PORT_SUCCESS).Return(nil)
	p.Receive(em)
	// em.AssertExpectations(t)

	_, err := em.Fields().ValueForPath("daemon")
	assert.NotNil(t, err, "field should not exists")
}
