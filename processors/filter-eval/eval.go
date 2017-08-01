//go:generate bitfanDoc
// Modify or add event's field with the result of an expression (math or compare)
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
	"fmt"
	"reflect"

	"github.com/Knetic/govaluate"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{
		opt:                 &options{},
		compiledExpressions: map[string]*govaluate.EvaluableExpression{},
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
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// If this filter is successful, remove arbitrary fields from this event. Example:
	// ` kv {
	// `   remove_field => [ "foo_%{somefield}" ]
	// ` }
	RemoveField []string `mapstructure:"remove_field"`

	// If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Example:
	// ` kv {
	// `   remove_tag => [ "foo_%{somefield}" ]
	// ` }
	// If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.
	RemoveTag []string `mapstructure:"remove_tag"`

	// list of field to set with expression's result
	// @ExampleLS expressions => { "usage" => "[usage] * 100" }
	Expressions map[string]interface{} `mapstructure:"expressions" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	var countError int
	for key, expressionString := range p.opt.Expressions {
		expression, err := p.cacheExpression(key, expressionString.(string))

		if err != nil {
			return fmt.Errorf("evaluation expression error : %s", err.Error())
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

		if err != nil {
			p.Logger.Errorf("error while evaluating `%s` with values `%s` ", expressionString, e.Fields())
			countError++
		} else {
			e.Fields().SetValueForPath(resultRaw, key)
		}
	}

	if countError == 0 {
		processors.ProcessCommonFields2(e.Fields(),
			p.opt.AddField,
			p.opt.AddTag,
			p.opt.RemoveField,
			p.opt.RemoveTag,
		)
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
	if e, ok := p.compiledExpressions[index]; ok {
		return e, nil
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

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expressionValue, functions)
	if err != nil {
		return nil, err
	}
	p.compiledExpressions[index] = expression

	return expression, nil
}
