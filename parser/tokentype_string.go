// generated by stringer -type=TokenType token.go; DO NOT EDIT

package parser

import "fmt"

const _TokenType_name = "LSTokenIllegalLSTokenEOFLSTokenUnexpectedEOFLSTokenWhitespaceLSTokenIdentifierLSTokenAssignmentLSTokenLCurlyBraceLSTokenLBracketLSTokenLParenLSTokenRParenLSTokenRCurlyBraceLSTokenRBracketLSTokenStringLSTokenNumberLSTokenIfLSTokenElseLSTokenElseIfLSTokenCommentLSTokenCommaLSTokenBool"

var _TokenType_index = [...]uint16{0, 14, 24, 44, 61, 78, 95, 113, 128, 141, 154, 172, 187, 200, 213, 222, 233, 246, 260, 272, 283}

func (i TokenType) String() string {
	if i < 0 || i+1 >= TokenType(len(_TokenType_index)) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
