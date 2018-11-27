//go:generate bitfanDoc
// Receive event on a ws connection
package websocketinput

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"bitfan/codecs"
	"bitfan/processors"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Receive event on a ws connection
func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The codec used for outputed data.
	// @Default "json"
	// @Type codec
	Codec codecs.CodecCollection

	// URI path
	// @Default "wsin"
	Uri string

	// Maximum message size allowed from peer.
	MaxMessageSize int `mapstructure:"max_message_size"`
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string

	wsupgrader websocket.Upgrader

	Hub *Hub
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Dec: codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Uri: "wsin",
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %v", err)
	}

	p.wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return err
}

func (p *processor) Start(e processors.IPacket) error {

	p.Hub = newHub(p.wellcome)
	go p.Hub.run()
	p.WebHook.Add(p.opt.Uri, p.HttpHandler)
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.Hub.stop()
	return nil
}

func (p *processor) processMessage(m []byte) {
	r := bytes.NewReader(m)

	var dec codecs.Decoder
	var err error
	if dec, err = p.opt.Codec.NewDecoder(r); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return
	}

	for dec.More() {
		var record interface{}
		if err = dec.Decode(&record); err != nil {
			if err == io.EOF {
				p.Logger.Debugln("error while decoding : ", err)
			} else {
				p.Logger.Errorln("error while decoding : ", err)
				break
			}
		}

		var e processors.IPacket
		switch v := record.(type) {
		case nil:
			continue
		case string:
			e = p.NewPacket(map[string]interface{}{"message": v})
		case map[string]interface{}:
			e = p.NewPacket(v)
		case []interface{}:
			e = p.NewPacket(map[string]interface{}{
				"request": v,
			})
		default:
			p.Logger.Errorf("Unknow structure %#v", v)
			continue
		}

		p.opt.ProcessCommonOptions(e.Fields())
		p.Send(e)
	}
}

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := p.wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		p.Logger.Errorf("websocket upgrade - %s", err.Error())
		return
	}

	client := &Client{
		hub:       p.Hub,
		conn:      conn,
		onMessage: p.processMessage,
	}
	p.Hub.register <- client

	go client.writePump()
	go client.readPump(int64(p.opt.MaxMessageSize))
}

func (p *processor) wellcome() [][]byte {
	return [][]byte{}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	onMessage func([]byte)

	done chan struct{}
}

// readPump pumps messages from the websocket connection to the hub.

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(maxMessageSize int64) {
	defer func() {
		c.conn.Close()
		close(c.done)
		c.hub.unregister <- c
	}()
	if maxMessageSize > 0 {
		c.conn.SetReadLimit(maxMessageSize)
	}

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("IsUnexpectedCloseError error: %v", err)
			}
			break
		}

		// Build PACKET
		c.onMessage(message)
	}

}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}

// hub maintains the set of active clients and Broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	done chan struct{}
}

func newHub(wellcomeMessage func() [][]byte) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		done:       make(chan struct{}),
	}
}

func (h *Hub) stop() {
	for c, _ := range h.clients {
		h.unregister <- c
	}
	close(h.done)
}

func (h *Hub) run() {
	for {
		select {
		case <-h.done:
			return
		case client := <-h.register:
			client.done = make(chan struct{}, 1)
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.conn.Close()
				delete(h.clients, client)
			}
		}
	}
}
