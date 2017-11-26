package logstash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenString(t *testing.T) {
	tok := token{Kind: 1, Value: "O"}
	assert.Equal(t, "TokenIllegal 'O'", tok.String())
}

func TestTokenStringValue(t *testing.T) {
	assert.Equal(t, "TokenIllegal", getTokenKindString(TokenIllegal))
	assert.Equal(t, "TokenEOF", getTokenKindString(TokenEOF))
	assert.Equal(t, "TokenAssignment", getTokenKindString(TokenAssignment))
	assert.Equal(t, "TokenLCurlyBrace", getTokenKindString(TokenLCurlyBrace))
	assert.Equal(t, "TokenLBracket", getTokenKindString(TokenLBracket))
	assert.Equal(t, "TokenRCurlyBrace", getTokenKindString(TokenRCurlyBrace))
	assert.Equal(t, "TokenRBracket", getTokenKindString(TokenRBracket))
	assert.Equal(t, "TokenString", getTokenKindString(TokenString))
	assert.Equal(t, "TokenNumber", getTokenKindString(TokenNumber))
	assert.Equal(t, "TokenIf", getTokenKindString(TokenIf))
	assert.Equal(t, "TokenElse", getTokenKindString(TokenElse))
	assert.Equal(t, "TokenElseIf", getTokenKindString(TokenElseIf))
	assert.Equal(t, "TokenComment", getTokenKindString(TokenComment))
	assert.Equal(t, "TokenComma", getTokenKindString(TokenComma))
	assert.Equal(t, "TokenBool", getTokenKindString(TokenBool))
	assert.Equal(t, "TokenIllegal", getTokenKindString(456))
}

func TestTokenStringHumanValue(t *testing.T) {
	assert.Equal(t, "TokenIllegal", getTokenKindHumanString(TokenIllegal))
	assert.Equal(t, "End of content", getTokenKindHumanString(TokenEOF))
	assert.Equal(t, "=>", getTokenKindHumanString(TokenAssignment))
	assert.Equal(t, "{", getTokenKindHumanString(TokenLCurlyBrace))
	assert.Equal(t, "[", getTokenKindHumanString(TokenLBracket))
	assert.Equal(t, "}", getTokenKindHumanString(TokenRCurlyBrace))
	assert.Equal(t, "]", getTokenKindHumanString(TokenRBracket))
	assert.Equal(t, "string", getTokenKindHumanString(TokenString))
	assert.Equal(t, "number", getTokenKindHumanString(TokenNumber))
	assert.Equal(t, "if", getTokenKindHumanString(TokenIf))
	assert.Equal(t, "else", getTokenKindHumanString(TokenElse))
	assert.Equal(t, "else if", getTokenKindHumanString(TokenElseIf))
	assert.Equal(t, "comment", getTokenKindHumanString(TokenComment))
	assert.Equal(t, ",", getTokenKindHumanString(TokenComma))
	assert.Equal(t, "bool", getTokenKindHumanString(TokenBool))
	assert.Equal(t, "TokenIllegal", getTokenKindHumanString(456))
}
