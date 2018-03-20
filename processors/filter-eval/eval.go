//go:generate bitfanDoc
// Modify or add event's field with the result of
//
// * an expression (math or compare)
// * a go template
//
// **Operators and types supported in expression :**
//
// * Modifiers: `+` `-` `/` `*` `&` `|` `^` `**` `%` `>>` `<<`
// * Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
// * Logical ops: `||` `&&`
// * Numeric constants, as 64-bit floating point (`12345.678`)
// * String constants (single quotes: `'foobar'`)
// * Date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
// * Boolean constants: `true` `false`
// * Parenthesis to control order of evaluation `(` `)`
// * Arrays (anything separated by `,` within parenthesis: `(1, 2, 'foo')`)
// * Prefixes: `!` `-` `~`
// * Ternary conditional: `?` `:`
// * Null coalescence: `??`
//
package evalprocessor

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/vjeantet/bitfan/commons"
	"github.com/vjeantet/bitfan/processors"
	"gopkg.in/Knetic/govaluate.v3"
)

func New() processors.Processor {
	return &processor{
		opt:                 &options{},
		compiledExpressions: map[string]*govaluate.EvaluableExpression{},
		compiledTemplates:   map[string]*template.Template{},
	}
}

const (
	PORT_SUCCESS = 0
)

func (p *processor) MaxConcurent() int { return 0 }

// Evaluate expression
type processor struct {
	processors.Base
	opt *options

	compiledExpressions map[string]*govaluate.EvaluableExpression
	compiledTemplates   map[string]*template.Template
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// list of field to set with expression's result
	// @ExampleLS expressions => { "usage" => "[usage] * 100" }
	Expressions map[string]interface{} `mapstructure:"expressions"`

	// list of field to set with a go template location
	// @ExampleLS templates => { "count" => "{{len .data}}", "mail"=>"mytemplate.tpl" }
	Templates map[string]string `mapstructure:"templates"`

	// You can set variable to be used in template by using ${var}.
	// each reference will be replaced by the value of the variable found in Template's content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{}
	p.opt = &defaults
	err = p.ConfigureAndValidate(ctx, conf, p.opt)

	if len(p.opt.Expressions)+len(p.opt.Templates) == 0 {
		return fmt.Errorf("set one expression or go template")
	}

	//Prepare templates
	for key, tplLocStr := range p.opt.Templates {
		loc, err := commons.NewLocation(tplLocStr, p.ConfigWorkingLocation)
		if err != nil {
			return err
		}
		tpl, _, err := loc.TemplateWithOptions(p.opt.Var)
		if err != nil {
			return err
		}
		p.compiledTemplates[key] = tpl
	}

	// Prepare expressions
	for key, expressionString := range p.opt.Expressions {
		if _, err = p.cacheExpression(key, expressionString.(string)); err != nil {
			return err
		}
	}

	return err
}

func (p *processor) Receive(e processors.IPacket) (err error) {

	var countError int

	if len(p.opt.Var) > 0 {
		e.Fields().SetValueForPath(p.opt.Var, "var")
	}

	// go templates
	for key, _ := range p.opt.Templates {
		buff := bytes.NewBufferString("")
		err = p.compiledTemplates[key].Execute(buff, e.Fields())
		if err != nil {
			p.Logger.Errorf("template %s error : %v", key, err)
			return err
		}
		e.Fields().SetValueForPath(buff.String(), key)
	}

	// govaluate expressions
	for key, expression := range p.compiledExpressions {
		parameters := EvaluatedParameters{}
		for _, v := range expression.Tokens() {
			if v.Kind == govaluate.VARIABLE {
				path := v.Value.(string)
				splits := strings.Split(path, ".")
				arrIndex := 0
				if len(splits) > 1 {
					arrIndex, err = strconv.Atoi(splits[len(splits)-1])
					if err == nil {
						// array index ! rebuild path
						path = strings.Join(splits[0:len(splits)-1], ".")
					}
				}

				paramValue, err := e.Fields().ValueForPath(path)

				if len(splits) > 1 {
					switch v := paramValue.(type) {
					case []interface{}:
						paramValue = v[arrIndex]
					case []string:
						paramValue = v[arrIndex]
					}
				}

				p.Logger.Debugf("paramValue-->%s", paramValue)
				if err != nil {
					continue
				}
				parameters[v.Value.(string)] = paramValue
			}
		}
		resultRaw, err := expression.Eval(parameters)

		if err != nil {
			p.Logger.Errorf("error while evaluating `%s` with values `%s` ", expression.String(), e.Fields())
			countError++
		} else {
			e.Fields().SetValueForPath(resultRaw, key)
		}
	}

	if countError == 0 {
		p.opt.ProcessCommonOptions(e.Fields())
	}

	p.Send(e)
	return nil
}

type EvaluatedParameters map[string]interface{}

func (ep EvaluatedParameters) Get(key string) (interface{}, error) {
	value, found := ep[key]

	// return `false` if parameter does not exist or value is nil
	if !found || value == nil {
		return false, nil
	}

	// return []interface{} if parameter is a slice of non-interface type
	// this is needed for govaluate to work with slices
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		interfaceSlice := make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			interfaceSlice[i] = s.Index(i).Interface()
		}
		return interfaceSlice, nil
	}

	return value, nil
}

func (p *processor) cacheExpression(index string, expressionValue string) (*govaluate.EvaluableExpression, error) {
	functions := map[string]govaluate.ExpressionFunction{

		"bool": func(args ...interface{}) (interface{}, error) {
			if len(args) == 0 {
				return false, nil
			}
			switch args[0].(type) {
			case bool:
				return args[0].(bool), nil
			default:
				return true, nil
			}
		},

		"len": func(args ...interface{}) (interface{}, error) {
			if len(args) == 0 {
				return float64(0), nil
			}
			switch reflect.TypeOf(args[0]).Kind() {

			case reflect.Slice:
				v := reflect.ValueOf(args[0])
				return (float64)(v.Len()), nil

			case reflect.String:
				length := len(args[0].(string))
				return (float64)(length), nil

			default:
				return (float64)(1), nil
			}
		},
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionValue, functions)
	if err != nil {
		return nil, err
	}
	p.compiledExpressions[index] = expression

	return expression, nil
}
