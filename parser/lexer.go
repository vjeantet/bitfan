package parser

import "github.com/vjeantet/go-lexer"

func lexBegin(l *lexer.L) lexer.StateFunc {
	l.SkipWhitespace()
	consumeComments(l)

	return lexIdent
}

// lexComment scans a comment. The left comment marker is known to be present.
func consumeComments(l *lexer.L) {
	var r rune
	r = l.Next()
	if r == '#' {
		l.Ignore()
		for {
			r = l.Next()
			if r == lexer.EOFRune {
				l.Emit(TokenComment)
				return
			}
			if r == '\n' || r == '\t' {
				l.Rewind()
				l.Emit(TokenComment)
				l.SkipWhitespace()
				if l.Next() == '#' {
					l.Rewind()
					consumeComments(l)
				}
				l.Rewind()
				break
			}
		}
	} else {
		l.Rewind()
	}
}

func lexIdent(l *lexer.L) lexer.StateFunc {

	if l.Next() == '"' {
		// l.Rewind()
		return lexString
	}

	l.Take("abcdefghijklmnopqrstuvwxyz._-1234567890\"'%")

	// Si il ne commence par par un "
	// 	Si il y a un { dans le current
	// 		prendre uniquement ce qu'il y a avant le premier {

	if l.Current() == "if" {
		return lexIf
	}

	if l.Current() == "else" {
		return lexElse
	}

	l.Emit(TokenIdentifier)

	var r rune
	// Des qu’il y a un IDENT, il y a un { ou ASSIGNEMENT
	l.SkipWhitespace()

	r = l.Next()
	if r == '{' {
		l.Rewind()
		return lexLeftDelim
	}
	if r == '=' {
		l.Rewind()
		return lexAssignement
	}
	l.Emit(TokenIllegal)
	return nil
}

func lexElse(l *lexer.L) lexer.StateFunc {
	l.Next()
	l.Next()
	l.Next()

	if l.Current() == "else if" {
		return lexElseIf
	}

	l.Rewind()
	l.Rewind()

	l.Emit(TokenElse)
	l.SkipWhitespace()

	r := l.Next()
	if r == 'i' && l.Next() == 'f' {
		return lexIf
	}

	if r == '{' {
		l.Rewind()
		return lexLeftDelim
	}
	l.Emit(TokenIllegal)
	return nil
}
func lexElseIf(l *lexer.L) lexer.StateFunc {
	l.Ignore()
	l.SkipWhitespace()

	var r rune
	for {
		r = l.Next()
		if r == ' ' && l.Peek() == '{' {
			l.Rewind()
			l.Emit(TokenElseIf)
			l.Next()
			l.Ignore()
			break
		}
		if r == '{' {
			l.Rewind()
			l.Emit(TokenElseIf)
			break
		}
	}

	return lexLeftDelim
}
func lexIf(l *lexer.L) lexer.StateFunc {
	l.Ignore()
	l.SkipWhitespace()

	var r rune
	for {
		r = l.Next()
		if r == ' ' && l.Peek() == '{' {
			l.Rewind()
			l.Emit(TokenIf)
			l.Next()
			l.Ignore()
			break
		}

		if r == '{' {
			l.Rewind()
			l.Emit(TokenIf)
			break
		}

	}

	return lexLeftDelim
}

func lexAssignement(l *lexer.L) lexer.StateFunc {
	l.Take("=>")
	if l.Current() == "=>" {
		l.Emit(TokenAssignment)
	} else {
		l.Emit(TokenIllegal)
		return nil
	}

	// Des qu’il y a un ASSIGNEMENT il y a un STRING ou NUMBER ou ARRAY
	return lexVariable
}

func lexVariable(l *lexer.L) lexer.StateFunc {
	var r rune
	l.SkipWhitespace()

	r = l.Next()
	if r == lexer.EOFRune {
		l.Emit(TokenEOF)
		return nil
	}

	if r == '"' || r == '\'' {
		return lexString
	}

	if lexer.IsDigit(r) {
		return lexNumber
	}

	if r == '[' {
		l.Emit(TokenLBracket)
		return lexVariable
	}

	if r == '{' {
		l.Emit(TokenLCurlyBrace)
		l.SkipWhitespace()
		consumeComments(l)
		return lexIdent
	}

	if lexer.IsLetter(r) {
		l.Rewind()
		return lexString
	}

	l.Emit(TokenIllegal)
	return nil
}

func lexNumber(l *lexer.L) lexer.StateFunc {
	l.Take("1234567890.")

	l.Emit(TokenNumber)

	l.SkipWhitespace()
	r := l.Next()
	if r == ']' {
		l.Emit(TokenRBracket)
		l.SkipWhitespace()
		consumeComments(l)
		r = l.Next()
	}

	if r == ',' {
		l.Emit(TokenComma)
		return lexVariable
	}

	if r == '}' {
		return lexRightDelim
	}
	if r == lexer.EOFRune {
		l.Emit(TokenEOF)
		return nil
	}

	return lexIdent
}

func lexString(l *lexer.L) lexer.StateFunc {
	var quoted bool = false
	if l.Current() == "\"" {
		quoted = true
	}
	if l.Current() == "'" {
		quoted = true
	}
	var last rune
	for {
		r := l.Next()

		if quoted && (r == '"' || r == '\'') && last != '\\' {
			l.Emit(TokenString)
			break
		}
		if !quoted {
			if lexer.IsSpace(r) || lexer.IsNewLine(r) || r == '}' || r == '{' {
				l.Rewind()
				if l.Current() == "true" || l.Current() == "false" {
					l.Emit(TokenBool)
				} else {
					l.Emit(TokenString)
				}

				break
			}
		}
		last = r
	}

	l.SkipWhitespace()
	consumeComments(l)
	r := l.Next()
	if r == ']' {
		l.Emit(TokenRBracket)
		l.SkipWhitespace()
		consumeComments(l)
		r = l.Next()
	}

	if r == '=' && l.Peek() == '>' {
		l.Rewind()
		return lexAssignement
	}

	if r == ',' {
		l.Emit(TokenComma)
		return lexVariable
	}

	if r == '}' {
		return lexRightDelim
	}

	if r == '{' {
		return lexLeftDelim
	}

	if r == lexer.EOFRune {
		l.Emit(TokenEOF)
		return nil
	}
	return lexIdent
}

// lexRightDelim scans the right delimiter, which is known to be present.
func lexRightDelim(l *lexer.L) lexer.StateFunc {
	l.Take("}")
	l.Emit(TokenRCurlyBrace)

	// Des qu’il y a un } il y a un IDENT ou } ou IF ou ELSEIF ou ELSE
	l.SkipWhitespace()
	consumeComments(l)
	r := l.Next()
	if r == '}' {
		return lexRightDelim
	}
	if r == lexer.EOFRune {
		l.Emit(TokenEOF)
		return nil
	}
	return lexIdent
}

// lexLeftDelim scans the left delimiter, which is known to be present.
func lexLeftDelim(l *lexer.L) lexer.StateFunc {
	l.Take("{")
	l.Emit(TokenLCurlyBrace)

	// Dés qu’il y a un { il y a un IDENT ou } ou IF
	l.SkipWhitespace()
	consumeComments(l)
	r := l.Next()
	if r == '}' {
		return lexRightDelim
	}
	if r == lexer.EOFRune {
		l.Emit(TokenEOF)
		return nil
	}

	return lexIdent
}
