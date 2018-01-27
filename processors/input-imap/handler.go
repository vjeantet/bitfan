package imap_input

import (
	"encoding/json"

	"github.com/vjeantet/bitfan/processors"
)

type toJsonHandler struct {
	packetFactory processors.PacketBuilder
	send          processors.PacketSender
	packet        processors.IPacket
}

func (hnd *toJsonHandler) Deliver(email string) error {
	docJSON, _ := json.Marshal(getMsg(email))
	e := hnd.packetFactory(map[string]interface{}{
		"message": string(docJSON),
	})
	hnd.send(e, 0)
	return nil
}

func (hnd *toJsonHandler) Describe() string {
	return "To JSON Handler"
}

func newToJsonHandler(pFactory processors.PacketBuilder, packet processors.IPacket, sender processors.PacketSender) *toJsonHandler {
	return &toJsonHandler{packetFactory: pFactory, packet: packet, send: sender}
}
