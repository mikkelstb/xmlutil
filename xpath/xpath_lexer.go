package xpath

import (
	"bufio"
	"io"
	"unicode"
)

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL
	ELEMENT
	DELIMETER
	ATTRIBUTE
	CONDITION_START
	CONDITION_END
	EQUALS
	VALUE
)

var tokens = []string{
	EOF:             "EOF",
	ILLEGAL:         "ILL",
	DELIMETER:       "DEL",
	ELEMENT:         "ELM",
	ATTRIBUTE:       "ATR",
	CONDITION_START: "[",
	CONDITION_END:   "]",
	EQUALS:          "=",
	VALUE:           "VAL",
}

func (t TokenType) String() string {
	return tokens[t]
}

type Token struct {
	Tokentype TokenType
	Position  Position
	Value     string
}

type Position int

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    0,
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) LexAll() []Token {
	var tokens []Token

	for {
		next_token := l.LexNext()
		tokens = append(tokens, next_token)

		if next_token.Tokentype == EOF {
			break
		}
	}
	return tokens
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.

func (l *Lexer) LexNext() Token {

	for {
		r, _, err := l.reader.ReadRune()
		l.pos++

		if err != nil {
			if err == io.EOF {
				return Token{Position: l.pos, Tokentype: EOF, Value: ""}
			}
			panic(err)
		}

		switch r {

		case '/':
			lit := l.lexDelimiter()
			return Token{Position: l.pos, Tokentype: DELIMETER, Value: lit}

		case '[':
			return Token{Position: l.pos, Tokentype: CONDITION_START, Value: string(r)}

		case ']':
			return Token{Position: l.pos, Tokentype: CONDITION_END, Value: string(r)}

		case '=':
			return Token{Position: l.pos, Tokentype: EOF, Value: string(r)}

		case '@':
			lit := "@" + l.lexAttribute()
			return Token{Position: l.pos, Tokentype: ATTRIBUTE, Value: lit}

		default:
			if unicode.IsSpace(r) {
				return Token{Position: l.pos, Tokentype: ILLEGAL, Value: "<space>"}

			} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
				start_pos := l.pos
				l.backup()
				lit := l.lexElement()
				return Token{Position: start_pos, Tokentype: ELEMENT, Value: lit}
			}
		}
	}
}

func (l *Lexer) lexDelimiter() string {
	r, _, err := l.reader.ReadRune()
	l.pos++
	if err != nil {
		if err == io.EOF {
			return "/"
		}
	}
	if r == '/' {
		return "//"
	}
	l.backup()
	return "/"
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.pos--
}

func (l *Lexer) lexAttribute() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexElement() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos++
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}
