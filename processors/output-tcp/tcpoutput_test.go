package tcpoutput

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor struct")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}

func TestLine(t *testing.T) {
	srv := &testServer{Address: ":3000"}
	assert.NoError(t, srv.Start(), "tcp server must be started")

	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"host": "localhost",
		"port": 3000,
		"codec": codecs.CodecCollection{
			Enc: codecs.New("line", map[string]interface{}{
				"format": "{{.message}}",
			}, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
	}

	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message1", map[string]interface{}{"abc": "def1", "n": 123})))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message2", map[string]interface{}{"abc": "def2", "n": 456})))
	assert.Equal(t, "message1\n", srv.GetMessage())
	assert.Equal(t, "message2\n", srv.GetMessage())
	assert.NoError(t, p.Stop(nil))
	srv.Stop()
}

type testServer struct {
	Address string
	listener net.Listener
	stringChan chan string
}

func (t *testServer) Start() error {
	var err error
	if t.listener, err = net.Listen("tcp", t.Address); err != nil {
		return err
	}
	t.stringChan = make(chan string, 10)
	go func() error {
		defer t.listener.Close()
		for {
			conn, err := t.listener.Accept()
			if err != nil {
				fmt.Println("TCP Fail accept with err:", err)
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					t.stringChan <- string(buf[:n])
				}
			}(conn)
		}
	}()
	return nil
}

func (t *testServer) Stop() {
	t.listener.Close()
	close(t.stringChan)
}

func (t *testServer) GetMessage() string {
	select {
	case m := <-t.stringChan: return m
	default: return ""
	}
}
