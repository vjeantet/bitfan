//go:generate bitfanDoc
package when

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/awillis/bitfan/processors"
	"gopkg.in/Knetic/govaluate.v3"
)

type processor struct {
	processors.Base

	opt                 *options
	compiledExpressions *sync.Map
}

type options struct {
	Expressions map[int]string
}

type EvaluatedParameters map[string]interface{}

func New() processors.Processor {
	return &processor{
		compiledExpressions: new(sync.Map),
		opt:                 &options{},
	}
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	switch v := conf["expressions"].(type) {
	case map[string]interface{}:
		conf["expressions"] = map[int]string{}
		for k, e := range v {
			ki, _ := strconv.Atoi(k)
			conf["expressions"].(map[int]string)[ki] = e.(string)
		}
	}
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

// comparison operators
// equality: ==, !=, <, >, <=, >=
// regexp: =~, !~
// inclusion: in, not in

// boolean operators
// and, or, nand, xor

// unary operators
// !

// Expressions can be long and complex. Expressions can contain other expressions,
// you can negate expressions with !, and you can group them with parentheses (...).

// 127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"

func (p *processor) Receive(e processors.IPacket) error {
	for order := 0; order < len(p.opt.Expressions); order++ {
		expressionValue := p.opt.Expressions[order]
		result, err := p.assertExpressionWithFields(order, expressionValue, e)
		if err != nil {
			p.Logger.Warnf("When processor evaluation '%s' error : %v", expressionValue, err)
			continue
		}

		if result {
			p.Send(e, order)
			break
		}
	}
	return nil
}

func (p *processor) assertExpressionWithFields(index int, expressionValue string, e processors.IPacket) (bool, error) {
	expression, err := p.cacheExpression(index, expressionValue)
	if err != nil {
		return false, fmt.Errorf("conditional expression error : %v", err)
	}
	parameters := EvaluatedParameters{}
	for _, v := range expression.Tokens() {
		if v.Kind == govaluate.VARIABLE {
			paramValues, err := e.Fields().ValuesForPath(v.Value.(string))
			if err != nil {
				continue
			}

			var paramValue interface{}
			switch len(paramValues) {
			case 0:
				paramValue = false
			case 1:
				paramValue = paramValues[0]
			default:
				paramValue = paramValues
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

func (p *processor) cacheExpression(index int, expressionValue string) (*govaluate.EvaluableExpression, error) {

	if e, ok := p.compiledExpressions.Load(index); ok {
		return e.(*govaluate.EvaluableExpression), nil
	}

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

	p.compiledExpressions.Store(index, expression)

	return expression, nil
}

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
