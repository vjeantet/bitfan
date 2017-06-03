package multilinecodec

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type multilineDecoder struct {
	more    bool
	r       *bufio.Scanner
	options multilineDecoderOptions
	memory  string
}

type multilineDecoderOptions struct {
	Delimiter string
	Negate    bool
	Pattern   string
	What      string
}

func New(r io.Reader, opt map[string]interface{}) *multilineDecoder {
	d := &multilineDecoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: multilineDecoderOptions{
			Delimiter: "\n",
			Negate:    false,
			Pattern:   "\\s",
			What:      "dd",
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

func (p *multilineDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	for p.r.Scan() {
		p.more = true
		match, _ := regexp.MatchString(p.options.Pattern, p.r.Text())

		if p.options.Negate {
			match = !match
		}

		if p.options.What == "previous" {
			if match == true { // stick to previous
				p.memory += p.options.Delimiter + p.r.Text()
				continue
			} else if p.memory == "" {
				p.memory = p.r.Text()
				continue
			} else {
				data["message"] = p.memory
				p.memory = p.r.Text()
				return data, nil
			}
		}
		if p.options.What == "next" {
			if match == true { // stick to au previous
				p.memory += "\n" + p.r.Text()
				continue
			} else if p.memory == "" {
				data["message"] = p.r.Text()
				p.memory = ""
				return data, nil
			} else {
				p.memory += "\n" + p.r.Text()
				data["message"] = p.memory
				p.memory = ""
				return data, nil
			}
		}
	}
	p.more = false
	return data, io.EOF

	// if p.options.What == "next" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + p.r.Text()
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		data["message"] = p.r.Text()
	// 		p.memory = ""
	// 		return data, nil
	// 	} else {
	// 		p.memory += "\n" + p.r.Text()
	// 		data["message"] = p.memory
	// 		p.memory = ""
	// 		return data, nil
	// 	}
	// }

	// if p.options.What == "previous" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + p.r.Text()
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		p.memory = p.r.Text()
	// 		return nil, nil
	// 	} else {
	// 		data["message"] = p.memory
	// 		p.memory = p.r.Text()
	// 		return data, nil
	// 	}
	// }

	return data, nil
}

func (p *multilineDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	// var cr io.Reader

	// var err error
	// cr, err = charset.NewReaderLabel(p.options.Charset, r)
	// if err != nil {
	// 	return nil, err
	// }
	// bytes, err := ioutil.ReadAll(cr)
	// if err != nil {
	// 	return data, err
	// }

	// match, _ := regexp.MatchString(p.options.Pattern, string(bytes))

	// if p.options.What == "next" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + string(bytes)
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		data["message"] = string(bytes)
	// 		p.memory = ""
	// 		return data, nil
	// 	} else {
	// 		p.memory += "\n" + string(bytes)
	// 		data["message"] = p.memory
	// 		p.memory = ""
	// 		return data, nil
	// 	}
	// }

	// if p.options.What == "previous" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + (string(bytes))
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		p.memory = string(bytes)
	// 		return nil, nil
	// 	} else {
	// 		data["message"] = p.memory
	// 		p.memory = string(bytes)
	// 		return data, nil
	// 	}
	// }

	return data, nil
}

func (p *multilineDecoder) More() bool {
	return p.more
}
