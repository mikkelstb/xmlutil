package xpath

import (
	"bufio"
	"io"
	"unicode"
)

type Token int

const (
	EOF Token = iota
	ILLEGAL
	ELEMENT
	DELIMETER
	ATTRIBUTE
	CONDITION_START
	CONDITION_END
	EQUALS
	VALUE
)

var tokens = []string {
	EOF:		"EOF",
	ILLEGAL:	"ILL",
	DELIMETER:	"DEL",
	ELEMENT:	"ELM",
	ATTRIBUTE:	"ATR",
	CONDITION_START: "[",
	CONDITION_END: "]",
	EQUALS: "=",
	VALUE: "VAL",
}

func (t Token) String() string {
	return tokens[t]
}



type Position int

type Lexer struct {
	pos Position
	reader *bufio.Reader
}


func NewLexer(reader io.Reader) *Lexer {
	return &Lexer {
		pos:	0,
		reader:	bufio.NewReader(reader),
	}
}


func (l *Lexer) LexAll() []Token {
	var tokens map[Token]string
	var next_token Token
}


// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.

func (l *Lexer) LexNext() (Position, Token, string) {

	for {
		r, _, err := l.reader.ReadRune()
		l.pos++

		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}

		switch r {

		case '/':
			lit := l.lexDelimiter()
			return l.pos, DELIMETER, lit

		case '[':
			return l.pos, CONDITION_START, "["

		case ']':
			return l.pos, CONDITION_START, "]"

		case '=':
			return l.pos, EQUALS, "="


		case '@':
			start_pos := l.pos
			lit := "@" + l.lexAttribute()
			return start_pos, ATTRIBUTE, lit

		default:
			if unicode.IsSpace(r) {
				return l.pos, ILLEGAL, "<space>"

			} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
				start_pos := l.pos
				l.backup()
				lit := l.lexElement()
				return start_pos, ELEMENT, lit			
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