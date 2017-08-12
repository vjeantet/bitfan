package doc

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Codec struct {
	Name       string
	PkgName    string
	ImportPath string
	Doc        string
	DocShort   string
	Decoder    *Decoder
	Encoder    *Encoder
}

type Encoder struct {
	Doc     string
	Options *CodecOptions
}
type Decoder struct {
	Doc     string
	Options *CodecOptions
}

type CodecOptions struct {
	Doc     string
	Options []*CodecOption
}
type CodecOption struct {
	Name           string
	Alias          string
	Doc            string
	Required       bool
	Type           string
	DefaultValue   interface{}
	PossibleValues []string
	//LogstashExample
	ExampleLS string
}

type Processor struct {
	Name       string
	ImportPath string
	Doc        string
	DocShort   string
	Options    *ProcessorOptions
	Ports      []*ProcessorPort
}

type ProcessorPort struct {
	Default bool
	Name    string
	Number  int
	Doc     string
}

type ProcessorOptions struct {
	Doc     string
	Options []*ProcessorOption
}

type ProcessorOption struct {
	Name           string
	Alias          string
	Doc            string
	Required       bool
	Type           string
	DefaultValue   interface{}
	PossibleValues []string
	//LogstashExample
	ExampleLS string
}

func (p *Processor) GenExample(kind string) []byte {
	g := &Generator{
		buf: bytes.Buffer{},
	}

	g.Printf("%s{\n", p.Name)
	for _, o := range p.Options.Options {
		g.Printf("\t%s\n", o.GenExample("logstash"))
	}
	g.Printf("}\n")
	return g.buf.Bytes()
}

func (p *Processor) GenMarkdown(kind string) []byte {
	g := &Generator{
		buf: bytes.Buffer{},
	}
	w := bufio.NewWriter(&g.buf)

	g.Printf("# %s\n", strings.ToUpper(p.Name))
	g.Printf("%s\n\n", p.Doc)
	if p.Options == nil {
		return g.buf.Bytes()
	}
	g.Printf("## Synopsys\n%s\n\n", p.Options.Doc)

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Setting", "Type", "Required", "Default value"})

	for _, o := range p.Options.Options {
		if o.Type == "processors.CommonOptions" {
			continue
		}
		required := "false"
		if o.Required {
			required = "true"
		}

		table.Append([]string{
			o.getIdentifier(), o.Type, required, o.getDefaultValue(),
		})
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.Render()
	w.Flush()

	g.Printf("\n\n")

	g.Printf("## Details\n\n")
	for _, o := range p.Options.Options {
		if o.Type == "processors.CommonOptions" {
			continue
		}
		g.Printf("### %s\n", o.getIdentifier())
		if o.Required {
			g.Printf("* This is a required setting.\n")
		}
		g.Printf("* Value type is %s\n", o.Type)
		g.Printf("* Default value is `%s`\n", o.getDefaultValue())
		g.Printf("\n%s\n\n", o.Doc)
	}

	g.Printf("\n\n")
	g.Printf("## Configuration blueprint\n\n")

	g.Printf("```\n%s{\n", p.Name)
	for _, o := range p.Options.Options {
		if o.Type == "processors.CommonOptions" {
			continue
		}
		g.Printf("\t%s\n", o.GenExample("logstash"))
	}
	g.Printf("}\n```\n")

	return g.buf.Bytes()
}

func (p *ProcessorOption) getIdentifier() string {
	if p.Alias == "" {
		return p.Name
	}
	return p.Alias
}

func (p *ProcessorOption) getDefaultValue() string {
	defaultValue := ""

	if p.DefaultValue != nil {
		defaultValue = p.DefaultValue.(string)
	} else {

		switch p.Type {
		case "hash":
			defaultValue = "{}"
		case "array":
			defaultValue = "[]"
		case "string":
			defaultValue = "\"\""
		case "int":
			defaultValue = "0"
		case "int64":
			defaultValue = "0"
		case "int32":
			defaultValue = "0"
		case "time.Duration":
			defaultValue = ""
		default:
			defaultValue = "?"
		}
	}
	return defaultValue
}

func (p *ProcessorOption) GenExample(kind string) string {
	example := p.ExampleLS
	if example == "" {

		if p.DefaultValue != nil {
			example = fmt.Sprintf("%s => %s", strings.ToLower(p.getIdentifier()), p.DefaultValue)
		} else {
			dv := ""
			switch p.Type {
			case "hash":
				dv = "{}"
			case "array":
				dv = "[]"
			case "string":
				dv = "\"\""
			case "int":
				dv = "123"
			case "int64":
				dv = "123"
			case "int32":
				dv = "123"
			case "time.Duration":
				dv = "30"
			default:
				dv = p.Type
			}
			example = fmt.Sprintf("%s => %s", strings.ToLower(p.getIdentifier()), dv)
		}

	}
	return example

}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}
