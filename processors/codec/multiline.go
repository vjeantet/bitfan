package codec

import (
	"io"
	"io/ioutil"
	"regexp"

	"golang.org/x/net/html/charset"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type MultilineDecoder struct {
	more    bool
	r       io.Reader
	options multilineDecoderOptions
	memory  string
}

type multilineDecoderOptions struct {
	Charset string

	Pattern string
	What    string
}

func NewMultilineDecoder(r io.Reader, opt map[string]interface{}) Decoder {
	d := &MultilineDecoder{
		r:    r,
		more: true,
		options: multilineDecoderOptions{
			Charset: "utf-8",
			Pattern: "\\s",
			What:    "dd",
		},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}
	pp.Println("-->", d.options)
	return d
}

func (p *MultilineDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	var cr io.Reader

	var err error
	cr, err = charset.NewReaderLabel(p.options.Charset, r)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(cr)
	if err != nil {
		return data, err
	}

	match, _ := regexp.MatchString(p.options.Pattern, string(bytes))

	if p.options.What == "next" {
		if match == true { // coller au previous
			p.memory += "\n" + string(bytes)
			return nil, nil
		} else if p.memory == "" {
			data["message"] = string(bytes)
			p.memory = ""
			return data, nil
		} else {
			p.memory += "\n" + string(bytes)
			data["message"] = p.memory
			p.memory = ""
			return data, nil
		}
	}

	if p.options.What == "previous" {
		if match == true { // coller au previous
			p.memory += "\n" + (string(bytes))
			return nil, nil
		} else if p.memory == "" {
			p.memory = string(bytes)
			return nil, nil
		} else {
			data["message"] = p.memory
			p.memory = string(bytes)
			return data, nil
		}
	}

	return data, nil
}

func (p *MultilineDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	p.more = false

	bytes, err := ioutil.ReadAll(p.r)
	if err != nil {
		return data, err
	}

	match, _ := regexp.MatchString(p.options.Pattern, string(bytes))

	if p.options.What == "next" {
		if match == true { // coller au previous
			p.memory += "\n" + string(bytes)
			return nil, nil
		} else if p.memory == "" {
			data["message"] = string(bytes)
			p.memory = ""
			return data, nil
		} else {
			p.memory += "\n" + string(bytes)
			data["message"] = p.memory
			p.memory = ""
			return data, nil
		}
	}

	if p.options.What == "previous" {
		if match == true { // coller au previous
			p.memory += "\n" + (string(bytes))
			return nil, nil
		} else if p.memory == "" {
			p.memory = string(bytes)
			return nil, nil
		} else {
			data["message"] = p.memory
			p.memory = string(bytes)
			return data, nil
		}
	}

	data["message"] = string(bytes)
	return data, nil
}

func (p *MultilineDecoder) More() bool {
	return p.more
}
