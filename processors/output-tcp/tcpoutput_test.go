package tcpoutput

import (
	"fmt"
	"github.com/awillis/bitfan/codecs"
	"github.com/awillis/bitfan/processors/doc"
	"github.com/awillis/bitfan/processors/testutils"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
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
	server, client := net.Pipe()
	srv := &testServer{conn: server}
	srv.Start()
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
	p.conn = client

	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message", map[string]interface{}{"abc": "def1", "n": 123})))
	time.Sleep(time.Second * 1)
	assert.Equal(t, "message\n", srv.GetMessage())
	assert.NoError(t, p.Stop(nil))
	srv.Stop()
}

type testServer struct {
	conn       net.Conn
	stringChan chan string
}

func (t *testServer) Start() {
	t.stringChan = make(chan string, 10)
	go func() error {
		for {
			buf := make([]byte, 1024)
			n, err := t.conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if n == 0 {
				fmt.Println("n=0")
				continue
			}
			t.stringChan <- string(buf[:n])
		}
	}()
}

func (t *testServer) Stop() {
	t.conn.Close()
	close(t.stringChan)
}

func (t *testServer) GetMessage() string {
	select {
	case m := <-t.stringChan:
		return m
	default:
		return ""
	}
}
