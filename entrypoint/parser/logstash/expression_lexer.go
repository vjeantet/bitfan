package logstash

import (
	"strings"

	"github.com/vjeantet/go-lexer"
)

var l lexer.L

const (
	UNKNOWN lexer.TokenType = iota

	PREFIX
	NUMERIC
	BOOLEAN
	STRING
	PATTERN
	TIME
	VARIABLE
	FUNCTION
	SEPARATOR

	COMPARATOR
	LOGICALOP
	MODIFIER

	CLAUSE
	CLAUSE_CLOSE

	TERNARY
)

var MODIFIER_SYMBOLS = []string{
	"+",
	"-",
	"*",
	"/",
	"%",
	"**",
	"&",
	"|",
	"^",
	">>",
	"<<",
}

var COMPARATOR_SYMBOLS = []string{
	"==",
	"!=",
	">",
	">=",
	"<",
	"<=",
	"=~",
	"!~",
	"in",
}

var REGEX_SYMBOLS = []string{
	"=~",
	"!~",
}

var LOGICAL_SYMBOLS = []string{
	"and", "or", "&&", "||",
}

var PREFIX_SYMBOLS = []string{
	"-", "!", "not", "~",
}

var SEPARATOR_SYMBOLS = []string{
	",",
}

var FUNCTION_NAMES = []string{
	"len", "bool",
}

/*
	Returns true if this operator is contained by the given array of candidate symbols.
	False otherwise.
*/
func isInSlice(needle string, candidates []string) bool {
	for _, symbolType := range candidates {
		if needle == symbolType {
			return true
		}
	}
	return false
}

func isQuote(r rune) bool {
	if r == '\'' || r == '"' {
		return true
	}
	return false
}

func lexBegin(l *lexer.L) lexer.StateFunc {
	l.SkipWhitespace()
	return lexIdent
}

func lexVariable(l *lexer.L) lexer.StateFunc {
	l.Take("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_1234567890-[]@")

	if l.Current() == "[" {
		l.Emit(CLAUSE)
	} else {
		l.Emit(VARIABLE)
	}

	return lexIdent
}

func lexString(l *lexer.L) lexer.StateFunc {
	l.Ignore()
	for {
		r := l.Next()
		if isQuote(r) {
			l.Rewind()
			l.Emit(STRING)
			l.Next()
			l.Ignore()
			break
		}
	}
	return lexIdent
}

func lexPattern(l *lexer.L) lexer.StateFunc {
	l.Ignore()
	for {
		r := l.Next()
		if r == '/' {
			l.Rewind()
			l.Emit(PATTERN)
			l.Next()
			l.Ignore()
			break
		}
	}
	return lexIdent
}

// ==, !=, <, >, <=, >= :: equality
// =~, !~ :: regexp
// in, not in :: inclusion
// and, or, nand, xor :: boolean operations
// ! not :: inversion
// () :: clause
// [abc][def] :: variables
// /.*/ :: regex pattern

//    `'foo' IN ('foo', 'bar')`
func lexIdent(l *lexer.L) lexer.StateFunc {
	l.SkipWhitespace()
	l.Ignore()

	r := l.Next()

	if r == '[' {
		return lexVariable
	}

	if r == '(' {
		l.Emit(CLAUSE)
		return lexIdent
	}
	if r == ')' || r == ']' {
		l.Emit(CLAUSE_CLOSE)
		return lexIdent
	}
	if r == ',' {
		l.Emit(SEPARATOR)
		return lexIdent
	}
	if isQuote(r) {
		return lexString
	}

	if lexer.IsDigit(r) {
		l.Take("1234567890.")
		l.Emit(NUMERIC)
		return lexIdent
	}

	l.Take("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz=!<>~*+-|&^")

	if isInSlice(strings.ToLower(l.Current()), COMPARATOR_SYMBOLS) {
		cu := l.Current()
		l.Emit(COMPARATOR)
		// check for regex pattern /.*/
		if cu == "=~" || cu == "!~" {
			l.SkipWhitespace()
			l.Ignore()
			r = l.Next()
			if r == '/' {
				return lexPattern
			} else {
				l.Rewind()
			}
		}
		return lexIdent
	}

	if isInSlice(strings.ToLower(l.Current()), LOGICAL_SYMBOLS) {
		l.Emit(LOGICALOP)
		return lexIdent
	}

	if isInSlice(strings.ToLower(l.Current()), MODIFIER_SYMBOLS) {
		l.Emit(MODIFIER)
		return lexIdent
	}

	if isInSlice(strings.ToLower(l.Current()), PREFIX_SYMBOLS) {
		l.Emit(PREFIX)
		return lexIdent
	}

	if isInSlice(strings.ToLower(l.Current()), []string{"true", "false"}) {
		l.Emit(BOOLEAN)
		return lexIdent
	}

	if isInSlice(strings.ToLower(l.Current()), FUNCTION_NAMES) {
		l.Emit(FUNCTION)
		return lexIdent
	}

	if r == lexer.EOFRune {
		return nil
	}

	l.Emit(FUNCTION)
	return lexIdent
}
