package logstash

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func readToken(stream *lexerStream) (token, error) {
	var ret token
	var tokenValue interface{}
	var tokenString string
	var kind tokenKind
	var character rune
	var completed bool
	var err error
	kind = TokenEOF

	for stream.canRead() {
		ret.Pos = stream.position
		character = stream.readCharacter()

		if unicode.IsSpace(character) {
			continue
		}

		kind = TokenIllegal

		if character == '#' {
			tokenString, _ = readUntilFalse(stream, true, false, true, isNotNewLine)
			tokenValue = tokenString
			kind = TokenComment
			break
		}

		if character == '{' {
			tokenValue = "{"
			kind = TokenLCurlyBrace
			break
		}

		if character == '}' {
			tokenValue = "}"
			kind = TokenRCurlyBrace
			break
		}

		// Assignment, accept = and =>
		if character == '=' || character == ':' {
			tokenValue = string(character)
			character = stream.readCharacter()
			if character == '>' {
				tokenValue = tokenValue.(string) + ">"
			} else {
				stream.rewind(1)
			}

			kind = TokenAssignment
			break
		}

		if character == ']' {
			tokenValue = "]"
			kind = TokenRBracket
			break
		}

		if character == '[' {
			tokenValue = "["
			kind = TokenLBracket
			break
		}

		if character == ',' {

			tokenValue = ","
			kind = TokenComma
			break
		}

		// number
		if isNumeric(character) {

			tokenString = readTokenUntilFalse(stream, isNumeric)

			if strings.Contains(tokenString, ".") {
				tokenValue, err = strconv.ParseFloat(tokenString, 64)
			} else {
				tokenValue, err = strconv.ParseInt(tokenString, 10, 64)
			}

			if err != nil {
				errorMsg := fmt.Sprintf("Unable to parse numeric value '%v'\n", tokenString)
				ret.Kind = kind
				ret.Value = tokenValue
				ret.Pos = stream.position
				ret.Line = stream.line
				ret.Col = stream.position - stream.lastEOLPos
				return ret, errors.New(errorMsg)
			}
			kind = TokenNumber
			break
		}

		// text
		if unicode.IsLetter(character) {
			stream.rewind(1)

			tokenString, _ = readUntilFalse(stream, true, false, true, isString)
			tokenValue = tokenString
			kind = TokenString

			if tokenString == "if" {
				tokenString, _ = readUntilFalse(stream, true, false, true, isNotLeftBr)
				tokenValue = tokenString
				kind = TokenIf
				break
			}

			if tokenString == "else" {
				kind = TokenElse
				tokenValue = "true"
				if stream.readCharacter() == ' ' {
					if stream.readCharacter() == 'i' {
						if stream.readCharacter() == 'f' {
							tokenString, _ = readUntilFalse(stream, true, false, true, isNotLeftBr)
							tokenValue = tokenString
							kind = TokenElseIf
							break
						}
						stream.rewind(1)
					}
					stream.rewind(1)
				}
				stream.rewind(1)
				break
			}

			// boolean?
			if tokenValue == "true" {
				kind = TokenBool
				tokenValue = true
			} else {
				if tokenValue == "false" {
					kind = TokenBool
					tokenValue = false
				}
			}

			break
		}

		if !isNotQuote(character) {
			tokenValue, completed = readUntilFalse(stream, true, false, true, isNotQuoteS(character))

			if !completed {
				return token{}, errors.New("Unclosed string literal")
			}

			// advance the stream one position, since reading until false assumes the terminator is a real token
			stream.rewind(-1)

			kind = TokenString
			break
		}

		errorMessage := fmt.Sprintf("Invalid token: '%v'", character)
		ret.Kind = kind
		ret.Value = string(character)
		ret.Pos = stream.position
		ret.Line = stream.line
		ret.Col = stream.position - stream.lastEOLPos
		return ret, errors.New(errorMessage)
	}

	ret.Kind = kind
	ret.Value = tokenValue
	ret.Pos = stream.position
	ret.Line = stream.line
	ret.Col = stream.position - stream.lastEOLPos
	return ret, nil
}

func readTokenUntilFalse(stream *lexerStream, condition func(rune) bool) string {

	var ret string

	stream.rewind(1)
	ret, _ = readUntilFalse(stream, false, true, true, condition)
	return ret
}

/*
	Returns the string that was read until the given [condition] was false, or whitespace was broken.
	Returns false if the stream ended before whitespace was broken or condition was met.
*/
func readUntilFalse(stream *lexerStream, includeWhitespace bool, breakWhitespace bool, allowEscaping bool, condition func(rune) bool) (string, bool) {

	var tokenBuffer bytes.Buffer
	var character rune
	var conditioned bool

	conditioned = false

	for stream.canRead() {

		character = stream.readCharacter()

		// Use backslashes to escape anything
		if allowEscaping && character == '\\' {

			character = stream.readCharacter()
			tokenBuffer.WriteString(string(character))
			continue
		}

		if unicode.IsSpace(character) {

			if breakWhitespace && tokenBuffer.Len() > 0 {
				conditioned = true
				break
			}
			if !includeWhitespace {
				continue
			}
		}

		if condition(character) {
			tokenBuffer.WriteString(string(character))
		} else {
			conditioned = true
			stream.rewind(1)
			break
		}
	}

	return tokenBuffer.String(), conditioned
}

func isNumeric(character rune) bool {

	return unicode.IsDigit(character) || character == '.' || character == '-'
}

func isNotNewLine(character rune) bool {

	return character != '\n'
}

func isNotLeftBr(character rune) bool {

	return character != '{'
}

func isNotQuoteS(character rune) func(rune) bool {
	return func(c rune) bool {
		if !isNotQuote(c) && c == character {
			return false
		}
		return true
	}
}

func isNotQuote(character rune) bool {
	return character != '\'' && character != '"'
}

func isNotAlphanumeric(character rune) bool {

	return !(unicode.IsDigit(character) ||
		unicode.IsLetter(character) ||
		character == '(' ||
		character == ')' ||
		!isNotQuote(character))
}

func isString(character rune) bool {
	return unicode.IsLetter(character) ||
		unicode.IsDigit(character) ||
		character == '_'
}
