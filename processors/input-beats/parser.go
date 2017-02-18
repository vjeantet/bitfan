package beatsinput

// This code come from https://github.com/packetzoom/logzoom
// https://github.com/packetzoom/logzoom/blob/master/input/filebeat/parser.go

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/vjeantet/bitfan/processors"
)

const (
	ack         = "2A"
	maxKeyLen   = 100 * 1024 * 1024 // 100 mb
	maxValueLen = 250 * 1024 * 1024 // 250 mb
)

type Parser struct {
	Conn       net.Conn
	Recv       chan map[string]interface{}
	wlen, plen uint32
	buffer     io.Reader
	Logger     processors.Logger
}

func NewParser(c net.Conn, r chan map[string]interface{}, log processors.Logger) *Parser {
	return &Parser{
		Conn:   c,
		Recv:   r,
		Logger: log,
	}
}

// ack acknowledges that the payload was received successfully
func (p *Parser) ack(seq uint32) error {
	buffer := bytes.NewBuffer([]byte(ack))
	binary.Write(buffer, binary.BigEndian, seq)
	//log.Printf("Sending ACK with seq %d", seq)

	if _, err := p.Conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// readKV parses key value pairs from within the payload
func (p *Parser) readKV() ([]byte, []byte, error) {
	var klen, vlen uint32

	// Read key len
	binary.Read(p.buffer, binary.BigEndian, &klen)

	if klen > maxKeyLen {
		return nil, nil, fmt.Errorf("key exceeds max len %d, got %d bytes", maxKeyLen, klen)
	}

	// Read key
	key := make([]byte, klen)
	_, err := p.buffer.Read(key)
	if err != nil {
		return nil, nil, err
	}

	// Read value len
	binary.Read(p.buffer, binary.BigEndian, &vlen)
	if vlen > maxValueLen {
		return nil, nil, fmt.Errorf("value exceeds max len %d, got %d bytes", maxValueLen, vlen)
	}

	// Read value
	value := make([]byte, vlen)
	_, err = p.buffer.Read(value)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

// read parses the compressed data frame
func (p *Parser) read() (uint32, error) {
	var seq, count uint32
	var k, v []byte
	var err error

	r, err := zlib.NewReader(p.Conn)
	if err != nil {
		return seq, err
	}
	defer r.Close()

	// Decompress
	buff := new(bytes.Buffer)
	io.Copy(buff, r)
	p.buffer = buff

	b := make([]byte, 2)
	for i := uint32(0); i < p.wlen; i++ {
		n, err := buff.Read(b)
		if err == io.EOF {
			return seq, err
		}

		if n == 0 {
			continue
		}
		switch string(b) {
		case "2D": // window size
			binary.Read(buff, binary.BigEndian, &seq)
			binary.Read(buff, binary.BigEndian, &count)

			var fields map[string]interface{}
			fields = make(map[string]interface{})

			for j := uint32(0); j < count; j++ {
				if k, v, err = p.readKV(); err != nil {
					return seq, err
				}
				fields[string(k)] = string(v)
			}

			if val, ok := fields["@timestamp"]; !ok {
				fields["@timestamp"] = time.Now()
			} else {
				fields["@timestamp"], _ = time.Parse("2006-01-02T15:04:05Z07:00", val.(string))
			}

			// Send to the receiver which is a buffer. We block because...
			p.Recv <- fields
		case "2J": // JSON
			//log.Printf("Got JSON data")
			binary.Read(buff, binary.BigEndian, &seq)
			binary.Read(buff, binary.BigEndian, &count)
			jsonData := make([]byte, count)
			_, err := p.buffer.Read(jsonData)
			//log.Printf("Got message: %s", jsonData)

			if err != nil {
				return seq, err
			}

			var fields map[string]interface{}

			decoder := json.NewDecoder(strings.NewReader(string(jsonData)))
			decoder.UseNumber()
			err = decoder.Decode(&fields)

			if val, ok := fields["@timestamp"]; !ok {
				fields["@timestamp"] = time.Now()
			} else {
				fields["@timestamp"], _ = time.Parse("2006-01-02T15:04:05Z07:00", val.(string))
			}

			if err != nil {
				return seq, err
			}

			// Send to the receiver which is a buffer. We block because...
			p.Recv <- fields

		default:
			return seq, fmt.Errorf("unknown type: %s", b)
		}
	}

	return seq, nil
}

// Parse initialises the read loop and begins parsing the incoming request
func (p *Parser) Parse() {
	b := make([]byte, 2)

Read:
	for {
		n, err := p.Conn.Read(b)

		if err != nil || n == 0 {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				p.Logger.Debugf("[%s] error reading %v", p.Conn.RemoteAddr().String(), err)
				break Read
			}
			p.Logger.Errorf("[%s] error reading %v", p.Conn.RemoteAddr().String(), err)
			break Read
		}

		switch string(b) {
		case "2W": // window length
			binary.Read(p.Conn, binary.BigEndian, &p.wlen)
		case "2C": // frame length
			binary.Read(p.Conn, binary.BigEndian, &p.plen)
			var seq uint32
			seq, err := p.read()

			if err != nil {
				log.Printf("[%s] error parsing %v", p.Conn.RemoteAddr().String(), err)
				break Read
			}

			if err := p.ack(seq); err != nil {
				log.Printf("[%s] error acking %v", p.Conn.RemoteAddr().String(), err)
				break Read
			}
		case "2J":

		default:
			// This really shouldn't happen
			p.Logger.Errorf("[%s] Received unknown type (%s): %s", p.Conn.RemoteAddr().String(), b, err)
			break Read
		}
	}
	close(p.Recv)
}
