//go:generate bitfanDoc -codec line
package linecodec

import (
	"bufio"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/codecs/lib"
)

// doc decoder
type decoder struct {
	more    bool
	r       *bufio.Scanner
	options decoderOptions

	log lib.Logger
}

// doc decoderOptions
type decoderOptions struct {
	// Change the delimiter that separates lines
	// @Default "\\n"
	Delimiter string
}

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: decoderOptions{
			Delimiter: "\n",
		},
	}

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Find the index of the input of a newline followed by a
		// pound sign.
		if i := strings.Index(string(data), d.options.Delimiter); i >= 0 {
			return i + 1, data[0:i], nil
		}

		// If at end of file with data return the data
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	d.r.Split(split)

	return d
}

func (d *decoder) SetOptions(conf map[string]interface{}, logger lib.Logger, cwl string) error {
	d.log = logger

	if err := mapstructure.Decode(conf, &d.options); err != nil {
		return err
	}
	return nil
}

func (d *decoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if true == d.r.Scan() {
		d.more = true
		data["message"] = d.r.Text()
	} else {
		d.more = false
		return data, io.EOF
	}

	return data, nil
}

func (d *decoder) More() bool {
	return d.more
}
