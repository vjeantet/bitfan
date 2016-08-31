package parser

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/vjeantet/go-lexer"
)

type Parser struct {
	l    *lexer.L
	line int
	col  int
}

type Configuration struct {
	Sections map[string]*Section
}

type Section struct {
	Name    string
	Plugins map[int]*Plugin
}

type Plugin struct {
	Name     string
	Codecs   map[int]*Codec
	Settings map[int]*Setting
	When     map[int]*When // IF and ElseIF with order
}

type Codec struct {
	Name     string
	Settings map[int]*Setting
}

type When struct {
	Expression string          // condition
	Plugins    map[int]*Plugin // actions
}

type Setting struct {
	K string
	V interface{}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{l: lexer.New(r, lexBegin)}
}

func (p *Parser) Parse() (*Configuration, error) {
	var err error
	var tok lexer.Token

	config := &Configuration{
		Sections: map[string]*Section{},
	}

	p.l.Start()
	for {

		tok, err = p.getToken(TokenComment, TokenIdentifier, TokenEOF, TokenRCurlyBrace)
		if err != nil {
			return config, fmt.Errorf("parse error Parse %s", err)
		}

		// If Comment Donoe
		if tok.Type == TokenEOF {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			var section *Section
			section, err = p.parseSection(&tok)
			if err != nil {
				return config, err
			}
			config.Sections[section.Name] = section
		}

	}

	return config, nil
}

func (p *Parser) parseSection(tok *lexer.Token) (*Section, error) {
	section := &Section{}
	if tok.Value != "input" && tok.Value != "filter" && tok.Value != "output" {
		return section, fmt.Errorf("parse error, unexpected '%s', line %d col %d", tok.Value, tok.Line, tok.Col)
	}

	section.Name = tok.Value
	section.Plugins = make(map[int]*Plugin, 0)

	// si pas de { alors erreur
	var err error
	*tok, err = p.getToken(TokenLCurlyBrace)
	if err != nil {
		return section, fmt.Errorf("section parse error %s", err)
	}
	i := 0
	for {
		*tok, err = p.getToken(TokenComment, TokenIdentifier, TokenRCurlyBrace, TokenIf, TokenElse, TokenElseIf)
		if err != nil {
			log.Printf(" -sp- %s %s", TokenType(tok.Type).String(), err)
			return section, fmt.Errorf("parse section error %s", err)
		}

		if tok.Type == TokenRCurlyBrace {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			plugin, err := p.parsePlugin(tok)
			if err != nil {
				return section, err
			}
			section.Plugins[i] = plugin
			i = i + 1
			continue
		case TokenIf:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return section, err
			}
			section.Plugins[i] = plugin
			i = i + 1
			continue
		case TokenElse:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return section, err
			}
			plugin.When[0].Expression = "true"
			iWhen := len(section.Plugins[i-1].When)
			section.Plugins[i-1].When[iWhen] = plugin.When[0]
			continue
		case TokenElseIf:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return section, err
			}
			iWhen := len(section.Plugins[i-1].When)
			section.Plugins[i-1].When[iWhen] = plugin.When[0]
			continue
		}
	}
	return section, nil
}

func (p *Parser) parseWHEN(tok *lexer.Token) (*Plugin, error) {
	pluginWhen := &Plugin{}
	pluginWhen.Name = "when"
	pluginWhen.When = make(map[int]*When, 0)

	when := &When{
		Expression: tok.Value,
		Plugins:    map[int]*Plugin{},
	}

	// si pas de { alors erreur
	var err error
	*tok, err = p.getToken(TokenLCurlyBrace)
	if err != nil {
		return pluginWhen, fmt.Errorf("IF parse error %s", err)
	}
	i := 0
	for {
		*tok, err = p.getToken(TokenComment, TokenIdentifier, TokenRCurlyBrace, TokenIf, TokenElse, TokenElseIf)
		if err != nil {
			return pluginWhen, fmt.Errorf("parse IF error %s", err)
		}

		if tok.Type == TokenRCurlyBrace {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			plugin, err := p.parsePlugin(tok)
			if err != nil {
				return pluginWhen, err
			}
			when.Plugins[i] = plugin
			i = i + 1
			continue
		case TokenIf:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return pluginWhen, err
			}
			when.Plugins[i] = plugin
			i = i + 1
			continue
		case TokenElse:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return pluginWhen, err
			}
			plugin.When[0].Expression = "true"
			iWhen := len(when.Plugins[i-1].When)
			when.Plugins[i-1].When[iWhen] = plugin.When[0]
			continue
		case TokenElseIf:
			plugin, err := p.parseWHEN(tok)
			if err != nil {
				return pluginWhen, err
			}
			iWhen := len(when.Plugins[i-1].When)
			when.Plugins[i-1].When[iWhen] = plugin.When[0]
			continue
		}
	}

	id := len(pluginWhen.When)
	pluginWhen.When[id] = when
	return pluginWhen, nil
}

func (p *Parser) parsePlugin(tok *lexer.Token) (*Plugin, error) {
	var err error

	plugin := &Plugin{}
	plugin.Name = tok.Value
	plugin.Settings = map[int]*Setting{}
	plugin.Codecs = map[int]*Codec{}
	// log.Printf(" -pp- %s %s", TokenType(tok.Type).String(), tok.Value)
	*tok, err = p.getToken(TokenLCurlyBrace)
	if err != nil {

		return plugin, fmt.Errorf("Plugin parse error %s", err)
	}
	i := 0
	var advancedTok *lexer.Token
	for {
		if advancedTok == nil {
			*tok, err = p.getToken(TokenComment, TokenIdentifier, TokenRCurlyBrace)
			if err != nil {
				return plugin, fmt.Errorf("plugin parse error %s", err)
			}
		} else {
			tok = advancedTok
		}

		if tok.Type == TokenRCurlyBrace {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			/*
				if tok.Value == "codec" {
					codec, err := p.parseCodec(tok)
					if err != nil {
						return plugin, err
					}

					tok2, err2 := p.getToken(TokenLCurlyBrace)
					if err2 != nil {
						// c'est pas un  { donc faut le remettre dans la boucle
						// log.Printf(" -pcs- %s %s", TokenType(tok2.Type).String(), tok2.Value)
						advancedTok = &tok2
						plugin.Codecs[i] = codec
						i = i + 1
						continue
					}

					settings, err := p.parseCodecSettings(&tok2)
					codec.Settings = settings
					plugin.Codecs[i] = codec
					i = i + 1
					continue
				}
			*/

			// Token is not a codec
			setting, err := p.parseSetting(tok)
			if err != nil {
				return plugin, err
			}
			plugin.Settings[i] = setting
			i = i + 1
			continue
		}

	}

	return plugin, nil
}

func (p *Parser) parseCodecSettings(tok *lexer.Token) (map[int]*Setting, error) {
	var err error
	settings := make(map[int]*Setting, 0)

	i := 0
	for {
		*tok, err = p.getToken(TokenComment, TokenIdentifier, TokenRCurlyBrace)
		if err != nil {
			return settings, fmt.Errorf("codec settings parse error %s", err)
		}

		if tok.Type == TokenRCurlyBrace {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			setting, err := p.parseSetting(tok)
			if err != nil {
				return settings, err
			}
			settings[i] = setting
			i = i + 1
			continue
		}

	}
	return settings, nil
}

func (p *Parser) parseCodec(tok *lexer.Token) (*Codec, error) {
	var err error

	codec := &Codec{}
	codec.Settings = map[int]*Setting{}

	// log.Printf(" -pc- %s %s", TokenType(tok.Type).String(), tok.Value)

	*tok, err = p.getToken(TokenAssignment)
	if err != nil {
		return codec, fmt.Errorf("codec 1 parse error %s", err)
	}

	// log.Printf(" -pc- %s %s", TokenType(tok.Type).String(), tok.Value)

	*tok, err = p.getToken(TokenString)
	if err != nil {
		return codec, fmt.Errorf("codec 2 parse error %s", err)
	}
	codec.Name = tok.Value
	// log.Printf(" -pc- %s %s", TokenType(tok.Type).String(), tok.Value)

	return codec, nil
}

func (p *Parser) parseSetting(tok *lexer.Token) (*Setting, error) {
	setting := &Setting{}

	if strings.HasPrefix(tok.Value, "\"") {
		tok.Value = strings.Replace(tok.Value, "\\", "", -1)
		tok.Value = strings.TrimPrefix(tok.Value, "\"")
		tok.Value = strings.TrimSuffix(tok.Value, "\"")
	}
	setting.K = tok.Value
	// log.Printf(" -- %s %s", TokenType(tok.Type).String(), tok.Value)

	var err error
	*tok, err = p.getToken(TokenAssignment)
	if err != nil {
		return setting, fmt.Errorf("Setting 1 parse error %s", err)
	}

	*tok, err = p.getToken(TokenString, TokenNumber, TokenLBracket, TokenLCurlyBrace, TokenBool)
	if err != nil {
		return setting, fmt.Errorf("Setting 2 parse error %s", err)
	}

	switch tok.Type {
	case TokenBool:
		setting.V = p.parseBool(tok.Value)
	case TokenString:
		setting.V = p.parseString(tok.Value)
	case TokenNumber:
		setting.V, err = p.parseNumber(tok.Value)
	case TokenLBracket:
		setting.V, err = p.parseArray()
	case TokenLCurlyBrace:
		setting.V, err = p.parseHash()
	}

	return setting, nil

}

func (p *Parser) parseBool(txt string) interface{} {
	var v interface{}
	// var err error
	if txt == "true" {
		v = true
	} else {
		v = false
	}
	return v
}

func (p *Parser) parseNumber(txt string) (interface{}, error) {
	var v interface{}
	var err error
	if strings.Contains(txt, ".") {
		v, err = strconv.ParseFloat(txt, 64)
	} else {
		v, err = strconv.ParseInt(txt, 10, 64)
	}
	return v, err
}

func (p *Parser) parseString(txt string) string {
	var v string
	if strings.HasPrefix(txt, "\"") {
		v = strings.Replace(txt, "\\", "", -1)
		v = strings.TrimPrefix(v, "\"")
		v = strings.TrimSuffix(v, "\"")
	} else {
		v = txt
	}
	return v
}

func (p *Parser) parseHash() (map[string]interface{}, error) {

	hash := map[string]interface{}{}
	for {
		tok, err := p.getToken(TokenComment, TokenIdentifier, TokenRCurlyBrace, TokenString, TokenComma)
		if err != nil {
			log.Fatalf("ParseHash parse error %s", err)
			return nil, err
		}

		if tok.Type == TokenRCurlyBrace {
			break
		}

		switch tok.Type {
		case TokenComment:
			continue
		case TokenIdentifier:
			set, err := p.parseSetting(&tok)
			if err != nil {
				return hash, err
			}
			hash[set.K] = set.V
		case TokenString:
			set, err := p.parseSetting(&tok)
			if err != nil {
				return hash, err
			}
			hash[set.K] = set.V
		}

	}
	return hash, nil
}

func (p *Parser) parseArray() ([]interface{}, error) {
	var str interface{}

	vals := make([]interface{}, 0, 20)
	for {
		tok, err := p.getToken(TokenString, TokenNumber, TokenComma, TokenRBracket)
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenRBracket {
			break
		}

		switch tok.Type {
		case TokenComma:
			continue
		case TokenNumber:
			str, err = p.parseNumber(tok.Value)
			if err != nil {
				return nil, err
			}
		case TokenString:
			if strings.HasPrefix(tok.Value, "\"") {
				v := strings.Replace(tok.Value, "\\", "", -1)
				v = strings.TrimPrefix(v, "\"")
				v = strings.TrimSuffix(v, "\"")
				str = v
			} else {
				str = tok.Value
			}
		}

		vals = append(vals, str)
	}

	return vals, nil
}

func (p *Parser) getToken(types ...lexer.TokenType) (lexer.Token, error) {

	tok, done := p.l.NextToken()

	if done {
		return lexer.Token{}, fmt.Errorf("unexpected end of file ")
	}

	if tok.Type == TokenIllegal {
		// log.Printf(" -- %s %s", TokenType(tok.Type).String(), tok.Value)
		return lexer.Token{}, fmt.Errorf("Illegal token '%s' found line %d col %d ", tok.Value, tok.Line, tok.Col)
	}

	for _, t := range types {
		if tok.Type == t {
			return *tok, nil
		}
	}

	if len(types) == 1 {
		return *tok, fmt.Errorf("unexpected token '%s' expected '%s' on line %d col %d", tok.Value, TokenType(types[0]).String(), tok.Line, tok.Col)
	}

	list := make([]string, len(types))
	for i, t := range types {
		list[i] = TokenType(t).String()
	}

	return *tok, fmt.Errorf("unexpected token '%s' expected one of %s on line %d col %d", tok.Value, strings.Join(list, "|"), tok.Line, tok.Col)
}

func (p *Parser) DumpTokens() {
	p.l.Start()
	for {
		tok, done := p.l.NextToken()
		if done {
			break
		}
		color := "\033[93m"
		if tok.Type == TokenIf || tok.Type == TokenElseIf || tok.Type == TokenElse {
			color = "\033[1m\033[91m"
		}
		if tok.Type == TokenLBracket || tok.Type == TokenRBracket || tok.Type == TokenRCurlyBrace || tok.Type == TokenLCurlyBrace {
			color = "\033[90m"
		}
		log.Printf("%4d line %3d:%-2d %s%-20s\033[0m _\033[92m%s\033[0m_", tok.Pos, tok.Line, tok.Col, color, TokenType(tok.Type).String(), tok.Value)
	}
}
