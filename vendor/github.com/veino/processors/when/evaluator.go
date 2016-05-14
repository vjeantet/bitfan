package when

// With vjeantet/conditions
// func assertExpressionWithFields(expressionValue string, fields *mxj.Map) (bool, error) {
// 	p := conditions.NewParser(strings.NewReader(expressionValue))
// 	expression, err := p.Parse()
// 	if err != nil {
// 		return false, err
// 	}
// 	parameters := make(map[string]interface{}, 8)
// 	for _, v := range conditions.Variables(expression) {
// 		paramValue, err := fields.ValueForPath(v)
// 		if err != nil {
// 			return false, fmt.Errorf("conditional field not found : %s", err.Error())
// 		}
// 		parameters[v] = paramValue
// 	}

// 	result, err := conditions.Evaluate(expression, parameters)

// 	return result, err
// }
