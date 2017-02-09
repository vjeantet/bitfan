//go:generate bitfanDoc
// This processor, depending on the condition evaluation, will route message to
// one or more different pipelines and/or pass the message through the processor to the next one.
// Behavior :
//
// * WHEN Condition is evaluated to true THEN the message go to the pipelines set in Path
// * WHEN Condition is evaluated to true AND Fork set to true THEN the message go to the pipeline set in Path AND pass through.
// * WHEN Condition is evaluated to false THEN the message pass through.
// * WHEN Condition is evaluated to false AND Fork set to true THEN the message  pass through.
package route

import (
	"fmt"

	"reflect"

	"github.com/Knetic/govaluate"
	"github.com/vjeantet/bitfan/processors"
)

const (
	PORT_SUCCESS = 0
	PORT_TRUNK   = 1
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this processor is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this processor is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// set a condition to fork and route message
	// when false, message is routed to trunk
	// By default condition is evaluated to true
	Condition string `mapstructure:"condition"`

	// Fork mode disabled by default
	// @Default false
	// @ExampleLS fork => false
	Fork bool `mapstructure:"fork"`

	// If this processor is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this processor is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

	// Add a type field to all events handled by this processor
	Type string

	// Path to configuration to send the incomming message, it could be a local file or an url
	// can be relative path to the current configuration.
	// @ExampleLS path=> ["error.conf"]
	Path []string `mapstructure:"path" validate:"required"`

	// You can set variable references in the used configuration by using ${var}.
	// each reference will be replaced by the value of the variable found in this option
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

// route message to other pipelines
type processor struct {
	processors.Base

	opt *options

	compiledExpression *govaluate.EvaluableExpression
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt = &options{
		Fork: false,
	}
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	var result bool = true
	if p.opt.Condition != "" {
		var err error
		result, err = p.assertExpressionWithFields(p.opt.Condition, e)
		if err != nil {
			p.Logger.Errorf("When processor evaluation error : %s\n", err.Error())
			return err
		}
	}

	processors.ProcessCommonFields2(e.Fields(),
		p.opt.Add_field,
		p.opt.Add_tag,
		p.opt.Remove_field,
		p.opt.Remove_tag,
	)

	if result == true && p.opt.Fork == true {
		p.Send(e.Clone(), PORT_TRUNK)
		p.Send(e, PORT_SUCCESS)
	} else if result == true && p.opt.Fork == false {
		p.Send(e, PORT_SUCCESS)
	} else {
		p.Send(e, PORT_TRUNK)
	}

	return nil
}

func (p *processor) assertExpressionWithFields(expressionValue string, e processors.IPacket) (bool, error) {
	expression, err := p.cacheExpression(expressionValue)
	if err != nil {
		return false, fmt.Errorf("conditional expression error : %s", err.Error())
	}
	parameters := EvaluatedParameters{}
	for _, v := range expression.Tokens() {
		if v.Kind == govaluate.VARIABLE {
			paramValue, err := e.Fields().ValueForPath(v.Value.(string))
			if err != nil {
				continue
			}
			parameters[v.Value.(string)] = paramValue
		}
	}
	resultRaw, err := expression.Eval(parameters)

	var result bool
	switch resultRaw.(type) {
	case bool:
		result = resultRaw.(bool)
	default:
		result = true
	}

	return result, err
}

func (p *processor) cacheExpression(expressionValue string) (*govaluate.EvaluableExpression, error) {
	if p.compiledExpression != nil {
		return p.compiledExpression, nil
	}

	functions := map[string]govaluate.ExpressionFunction{

		"bool": func(args ...interface{}) (interface{}, error) {
			switch args[0].(type) {
			case bool:
				return args[0].(bool), nil
			default:
				return true, nil
			}
		},

		"len": func(args ...interface{}) (interface{}, error) {
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

	var err error
	p.compiledExpression, err = govaluate.NewEvaluableExpressionWithFunctions(expressionValue, functions)
	if err != nil {
		return nil, err
	}

	return p.compiledExpression, nil
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
