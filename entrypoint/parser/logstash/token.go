package logstash

import "fmt"

// Represents a single parsed token.
type token struct {
	Kind  tokenKind
	Value interface{}
	Pos   int
	Line  int
	Col   int
}

// Represents all valid types of tokens that a token can be.
type tokenKind int

const (
	TokenIllegal tokenKind = iota + 1
	TokenEOF
	TokenAssignment
	TokenLCurlyBrace
	TokenLBracket
	TokenRCurlyBrace
	TokenRBracket
	TokenString
	TokenNumber
	TokenIf
	TokenElse
	TokenElseIf
	TokenComment
	TokenComma
	TokenBool
)

func (t *token) String() string {
	return fmt.Sprintf("%s '%s'", getTokenKindString(t.Kind), t.Value)
}

// GetTokenKindString returns a string that describes the given tokenKind.
func getTokenKindString(kind tokenKind) string {

	switch kind {

	case TokenIllegal:
		return "TokenIllegal"
	case TokenEOF:
		return "TokenEOF"
	case TokenAssignment:
		return "TokenAssignment"
	case TokenLCurlyBrace:
		return "TokenLCurlyBrace"
	case TokenLBracket:
		return "TokenLBracket"
	case TokenRCurlyBrace:
		return "TokenRCurlyBrace"
	case TokenRBracket:
		return "TokenRBracket"
	case TokenString:
		return "TokenString"
	case TokenNumber:
		return "TokenNumber"
	case TokenIf:
		return "TokenIf"
	case TokenElse:
		return "TokenElse"
	case TokenElseIf:
		return "TokenElseIf"
	case TokenComment:
		return "TokenComment"
	case TokenComma:
		return "TokenComma"
	case TokenBool:
		return "TokenBool"
	}

	return "TokenIllegal"
}

func getTokenKindHumanString(kind tokenKind) string {

	switch kind {

	case TokenIllegal:
		return "TokenIllegal"
	case TokenEOF:
		return "End of content"
	case TokenAssignment:
		return "=>"
	case TokenLCurlyBrace:
		return "{"
	case TokenLBracket:
		return "["
	case TokenRCurlyBrace:
		return "}"
	case TokenRBracket:
		return "]"
	case TokenString:
		return "string"
	case TokenNumber:
		return "number"
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenElseIf:
		return "else if"
	case TokenComment:
		return "comment"
	case TokenComma:
		return ","
	case TokenBool:
		return "bool"
	}

	return "TokenIllegal"
}
