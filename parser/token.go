package parser

import "github.com/vjeantet/go-lexer"

//go:generate stringer -type=TokenType token.go

type TokenType lexer.TokenType

const (
	LSTokenIllegal TokenType = iota
	LSTokenEOF
	LSTokenUnexpectedEOF
	LSTokenWhitespace
	LSTokenIdentifier
	LSTokenAssignment
	LSTokenLCurlyBrace
	LSTokenLBracket
	LSTokenLParen
	LSTokenRParen
	LSTokenRCurlyBrace
	LSTokenRBracket
	LSTokenString
	LSTokenNumber
	LSTokenIf
	LSTokenElse
	LSTokenElseIf
	LSTokenComment
	LSTokenComma
	LSTokenBool
)

const (
	TokenIllegal lexer.TokenType = iota
	TokenEOF
	TokenUnexpectedEOF
	TokenWhitespace
	TokenIdentifier
	TokenAssignment
	TokenLCurlyBrace
	TokenLBracket
	TokenLParen
	TokenRParen
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
