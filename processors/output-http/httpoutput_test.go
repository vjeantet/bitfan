package httpoutput

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"bitfan/codecs"
	"bitfan/processors/doc"
	"bitfan/processors/testutils"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor struct")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}

func TestDefault(t *testing.T) {
	c := make(chan map[string]interface{}, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		var j map[string]interface{}
		if err := json.Unmarshal(body, &j); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		c <- j
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"url":            ts.URL,
		"batch_size":     1,
		"batch_interval": 1,
	}

	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("msg 1", map[string]interface{}{"abc1": "def1", "1": 123, "@timestamp": "ts"})))
	assert.Equal(t, map[string]interface{}{"message": "msg 1", "abc1": "def1", "1": 123.0, "@timestamp": "ts"}, <-c)
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message 2", map[string]interface{}{"abc2": "def2", "2": 456, "@timestamp": "ts"})))
	assert.Equal(t, map[string]interface{}{"message": "message 2", "abc2": "def2", "2": 456.0, "@timestamp": "ts"}, <-c)
	assert.NoError(t, p.Stop(nil))
}
func TestLine(t *testing.T) {
	c := make(chan string, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		c <- string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"url": ts.URL,
		"codec": codecs.CodecCollection{
			Enc: codecs.New("line", map[string]interface{}{
				"format": "{{.message}}\t{{.abc}}\t{{.n}}",
			}, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		"batch_size":     3,
		"batch_interval": 1,
	}

	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message1", map[string]interface{}{"abc": "def1", "n": 123})))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message2", map[string]interface{}{"abc": "def2", "n": 456})))
	assert.Equal(t, "message1\tdef1\t123\nmessage2\tdef2\t456\n", <-c)
	assert.NoError(t, p.Stop(nil))
}

func TestRetry(t *testing.T) {
	c := make(chan string, 1)
	counter := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		if counter == 1 {
			ioutil.ReadAll(r.Body)
			w.WriteHeader(http.StatusInternalServerError)
			c <- "500"
		} else {
			body, _ := ioutil.ReadAll(r.Body)
			c <- string(body)
			w.WriteHeader(http.StatusOK)
		}
		r.Body.Close()
	}))
	defer ts.Close()
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"url": ts.URL,
		"codec": codecs.CodecCollection{
			Enc: codecs.New("line", map[string]interface{}{
				"format": "{{.message}}",
			}, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		"batch_size":     1,
		"batch_interval": 1,
		"retry_interval": 1,
	}
	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("message", nil)))
	assert.Equal(t, "500", <-c)
	assert.Equal(t, "message\n", <-c)
	assert.NoError(t, p.Stop(nil))
}

func TestStopInRetry(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"url":            "bad url",
		"batch_size":     1,
		"batch_interval": 1,
		"retry_interval": 1,
	}
	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.NoError(t, p.Start(nil))
	assert.NoError(t, p.Receive(testutils.NewPacketOld("doom message", nil)))
	time.Sleep(time.Second)
	assert.NoError(t, p.Stop(nil))
}
