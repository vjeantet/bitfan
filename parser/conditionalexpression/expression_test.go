package conditionalexpression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressions(t *testing.T) {
	check(t, "true", `true`)
	check(t,
		`[testInt] == 4`,
		`[testInt] == 4`,
	)
	check(t,
		`[testInt] == len(4)`,
		`[testInt] == len ( 4 )`,
	)
	check(t,
		`[testInt][test] == 4`,
		`[testInt.test] == 4`,
	)
	check(t,
		`[testInt] == 8/2`,
		`[testInt] == 8 / 2`,
	)
	check(t,
		`[testInt] == [testInt3]+1`,
		`[testInt] == [testInt3] + 1`,
	)
	check(t,
		`!(false)`,
		`! ( false )`,
	)
	check(t,
		`"_grokparsefailure" in [tags]`,
		`'_grokparsefailure' in [tags]`,
	)
	check(t,
		`"_mumu" not in [tags]`,
		`! ( '_mumu' in [tags] )`,
	)
	check(t,
		`"_mumu" not in [tags] and [way] == "SEND"`,
		`! ( '_mumu' in [tags] && [way] == 'SEND' )`,
	)
	check(t,
		`[testString] == "true"`,
		`[testString] == 'true'`,
	)
	check(t,
		`[location][city] == "Paris"`,
		`[location.city] == 'Paris'`,
	)
	check(t,
		`[testInt] == 4`,
		`[testInt] == 4`,
	)
	check(t,
		`[way]`,
		`[way]`,
	)
	check(t,
		`[testInt] > [testInt3]`,
		`[testInt] > [testInt3]`,
	)
	check(t,
		`true and true`,
		`true && true`,
	)
	check(t,
		`true && true`,
		`true && true`,
	)
	check(t,
		`true or false`,
		`true || false`,
	)
	check(t,
		`true || false`,
		`true || false`,
	)
	check(t,
		`"foo" in ("foor", "foo","bar")`,
		`'foo' in ( 'foor' , 'foo' , 'bar' )`,
	)
	check(t,
		`"foo" not in ("foos", "sfoo","bar")`,
		`! ( 'foo' in ( 'foos' , 'sfoo' , 'bar' ) )`,
	)
	check(t,
		`[way] =~ '(RECEIVE|SEND)'`,
		`[way] =~ '(RECEIVE|SEND)'`,
	)
	check(t,
		`[way] =~ /(RECEIVE|SEND)/`,
		`[way] =~ '(RECEIVE|SEND)'`,
	)

	check(t,
		`not(true)`,
		`! ( true )`,
	)
	check(t,
		`"grokparsefailure" in [tags]`,
		`'grokparsefailure' in [tags]`,
	)
	check(t,
		`"_mumu" in [tags]`,
		`'_mumu' in [tags]`,
	)
	check(t,
		`"_grokparsefailure" not in [tags] || [way] != "SEND"`,
		`! ( '_grokparsefailure' in [tags] ) || [way] != 'SEND'`,
	)
	check(t,
		`[testString] == "false"`,
		`[testString] == 'false'`,
	)
	check(t,
		`[testString] != "true"`,
		`[testString] != 'true'`,
	)
	check(t,
		`[location][city] != "Paris"`,
		`[location.city] != 'Paris'`,
	)
	check(t,
		`[testInt] > 30`,
		`[testInt] > 30`,
	)
	check(t,
		`true and false`,
		`true && false`,
	)
	check(t,
		`true && false`,
		`true && false`,
	)
	check(t,
		`false or false`,
		`false || false`,
	)
	check(t,
		`false || false`,
		`false || false`,
	)
	check(t,
		`"foo" in ("foor", "foos","bar")`,
		`'foo' in ( 'foor' , 'foos' , 'bar' )`,
	)
	check(t,
		`"foo" not in ("foo", "sfoo","bar")`,
		`! ( 'foo' in ( 'foo' , 'sfoo' , 'bar' ) )`,
	)
	check(t,
		`[testUnk] == 3`,
		`[testUnk] == 3`,
	)
	check(t,
		`[testUnk]`,
		`[testUnk]`,
	)
	check(t,
		`[way] !~ '(RECEIVE|SEND)'`,
		`[way] !~ '(RECEIVE|SEND)'`,
	)
	check(t,
		`[way] !~ /(RECEIVE|SEND)/`,
		`[way] !~ '(RECEIVE|SEND)'`,
	)
}

func check(t *testing.T, lsExpression string, gvExpression string) {
	result, err := ToWhenExpression(lsExpression)
	assert.NoError(t, err, "err is not nil")
	assert.Equal(t, result, gvExpression)
}
