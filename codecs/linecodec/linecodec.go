package linecodec

import (
	"bufio"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type lineDecoder struct {
	more    bool
	r       *bufio.Scanner
	options options
}

type options struct {
	Charset string

	// Change the delimiter that separates lines
	// @Default "\\n"
	Delimiter string
}

func New(r io.Reader, opt map[string]interface{}) *lineDecoder {

	d := &lineDecoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: options{
			Charset:   "utf-8",
			Delimiter: "\n",
		},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
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

func (c *lineDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	if true == c.r.Scan() {
		c.more = true
		data["message"] = c.r.Text()
	} else {
		c.more = false
		return data, io.EOF
	}

	return data, nil
}

func (c *lineDecoder) More() bool {
	return c.more
}

func (c *lineDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	// var cr io.Reader

	// var err error
	// cr, err = charset.NewReaderLabel(c.options.Charset, r)
	// if err != nil {
	// 	return nil, err
	// }
	// csvr := csv.NewReader(cr)
	// csvr.Comma = c.comma

	// record, err := csvr.Read()

	// if err == io.EOF {
	// 	return data, err
	// }
	// if c.columnnames == nil {
	// 	c.columnnames = record
	// 	return nil, nil
	// }

	// for i, v := range c.columnnames {
	// 	data[v] = record[i]
	// 	// data[fmt.Sprintf("col_%d", i)] = v
	// }

	return data, nil
}
