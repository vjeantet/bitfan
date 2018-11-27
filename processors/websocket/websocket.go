//go:generate bitfanDoc
// Send event received over a ws connection to connected clients
package websocket

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"bitfan/codecs"
	"bitfan/processors"
)

// expose events thought a websocket
func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// The codec used for outputed data.
	// @Default "json"
	// @Type codec
	Codec codecs.CodecCollection

	// URI path
	// @Default "wsout"
	Uri string
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string

	Hub *Hub
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Enc: codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Uri: "wsout",
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %v", err)
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

func (p *processor) wellcome() [][]byte {
	r := [][]byte{}
	if m, ok := p.Memory.Get("last"); ok {
		r = append(r, m.([]byte))
	}
	return r
}

func (p *processor) Receive(e processors.IPacket) error {

	w := bytes.NewBuffer(nil)

	// Encode content
	var err error
	var enc codecs.Encoder
	enc, err = p.opt.Codec.NewEncoder(w)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())
		return err
	}

	enc.Encode(e.Fields().Old())

	p.Memory.Set("last", w.Bytes())
	for c, _ := range p.Hub.clients {
		c.send <- w.Bytes()
	}
	return nil
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: p.Hub, conn: conn, send: make(chan []byte, 256)}
	p.Hub.register <- client
	for _, m := range p.Hub.Wellcome() {
		client.send <- m
	}

	go client.writePump()
	go client.readPump()
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("IsUnexpectedCloseError error: %v", err)
			}
			break
		}
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.Broadcast <- message
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
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	// w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// hub maintains the set of active clients and Broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	Wellcome func() [][]byte

	done chan struct{}
}

func newHub(wellcomeMessage func() [][]byte) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Wellcome:   wellcomeMessage,
		done:       make(chan struct{}),
	}
}

func (h *Hub) stop() {
	for c, _ := range h.clients {
		h.unregister <- c
		c.conn.Close()
	}
	close(h.done)
}

func (h *Hub) run() {
	for {
		select {
		case <-h.done:
			return
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
