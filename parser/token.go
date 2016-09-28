package parser

// Represents a single parsed token.
type Token struct {
	Kind  TokenKind
	Value interface{}
	Pos   int
	Line  int
	Col   int
}

// Represents all valid types of tokens that a token can be.
type TokenKind int

const (
	TokenIllegal TokenKind = iota + 1
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

// GetTokenKindString returns a string that describes the given TokenKind.
func GetTokenKindString(kind TokenKind) string {

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
