package parser

type lexerStream struct {
	source     []rune
	position   int
	length     int
	line       int
	lastEOLPos int
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
	ret.line = 1
	return ret
}

func (this *lexerStream) readCharacter() rune {

	var character rune

	character = this.source[this.position]
	this.position += 1
	if character == '\n' {
		this.line += 1
		this.lastEOLPos = this.position
	}
	return character
}

func (this *lexerStream) rewind(amount int) {
	this.position -= amount

	if amount > 0 {
		if this.source[this.position] == '\n' {
			this.line -= 1
			this.lastEOLPos = this.position
		}
	}
}

func (this lexerStream) canRead() bool {
	return this.position < this.length
}
