package logstash

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/vjeantet/go-lexer"
)

// returns a bitfan compatible expression from a logstash one
func toWhenExpression(lsExpression string) (string, error) {
	lsExpression = fixNotInExpr(lsExpression)
	r := bytes.NewReader([]byte(lsExpression))
	l := lexer.New(r, lexBegin)
	l.Start()

	ktokens := ""
	separator := ""
	for {
		tok, done := l.NextToken()
		if done {
			break
		}

		value := tok.Value

		switch tok.Type {
		case PREFIX:
			if value == "not" {
				value = "!"
			}
		case NUMERIC:
		case BOOLEAN:
		case STRING:
			value = "'" + value + "'"
		case PATTERN:
			value = "'" + value + "'"
		case TIME:
			value = "'" + value + "'"
		case VARIABLE:
			value = strings.Replace(value, `][`, `.`, -1)

		case FUNCTION:
		case SEPARATOR:
		case COMPARATOR:
		case LOGICALOP:
			if value == "and" {
				value = "&&"
			}
			if value == "or" {
				value = "||"
			}
		case MODIFIER:
		case CLAUSE:
			value = "("
		case CLAUSE_CLOSE:
			value = ")"

		case TERNARY:

		default:
			return ktokens, fmt.Errorf("unknow token %s", tok.Type)
		}
		ktokens = ktokens + separator + value
		separator = " "
	}

	return ktokens, nil
}

// fixNotInExpr converts `x not in y` into `! (x in y)` as govaluate doesn't have `not in` operators
func fixNotInExpr(exprIn string) (exprOut string) {
	if !strings.Contains(exprIn, "not in") {
		return exprIn
	}

	re := regexp.MustCompile(` (&&|\|\|) `)
	andOrOperators := re.FindAllString(exprIn, -1)
	splitExpr := re.Split(exprIn, -1)

	for i, e := range splitExpr {
		if strings.Contains(e, "not in") {
			e = strings.Replace(e, "not in", "in", 1)
			splitExpr[i] = fmt.Sprintf("! (%s)", e)
		}
	}

	for i, e := range splitExpr {
		exprOut += e
		if len(andOrOperators) > i {
			exprOut += andOrOperators[i]
		}
	}
	return exprOut
}
