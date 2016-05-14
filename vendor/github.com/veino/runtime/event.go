package runtime

import (
	"time"

	"github.com/clbanning/mxj"
	"github.com/veino/veino"
)

// Events types constants
const (
	_ = iota
	genericEvent
	startEvent
	stopEvent
	tickEvent
)

// event represents data sent to agents (or received by agents)
type event struct {
	kind   int
	fields mxj.Map
}

func (e *event) Kind() int {
	return e.kind
}

func (e *event) SetKind(k int) {
	e.kind = k
}

func (e *event) Fields() *mxj.Map {
	return &e.fields
}

func (e *event) SetFields(f mxj.Map) {
	e.fields = f
}

func (e *event) Message() string {
	return e.Fields().ValueOrEmptyForPathString("message")
}

func (e *event) SetMessage(s string) {
	e.Fields().SetValueForPath(s, "message")
}

func (e *event) Clone() veino.IPacket {
	nf, _ := e.Fields().Copy()
	return NewPacket(e.Message(), nf)
}

func NewPacket(message string, fields map[string]interface{}) veino.IPacket {
	if fields == nil {
		fields = mxj.Map{}
	}
	fields["message"] = message

	if _, k := fields["@timestamp"]; !k {
		fields["@timestamp"] = time.Now().Format(veino.VeinoTime)
	}
	return &event{
		kind:   genericEvent,
		fields: fields,
	}
}
