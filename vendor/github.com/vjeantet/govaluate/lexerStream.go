package govaluate

type lexerStream struct {
	source   []rune
	position int
	length   int
}

func newLexerStream(source string) *lexerStream {

	var ret *lexerStream
	var runes []rune

	for _, character := range source {
		runes = append(runes, character)
	}

	ret = new(lexerStream)
	ret.source = runes
	ret.length = len(runes)
	return ret
}

func (l *lexerStream) readCharacter() rune {

	var character rune

	character = l.source[l.position]
	l.position++
	return character
}

func (l *lexerStream) rewind(amount int) {
	l.position -= amount
}

func (l lexerStream) canRead() bool {
	return l.position < l.length
}
