package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func startState(l *Lexer) stFunc {

	r := l.next()

	switch r {
	case '/':
		return delimeterState
	case '@':
		return attributeState
	case '[':
		return conditionStartState
	case ']':
		return conditionEndState
	}

	if strings.ContainsAny(string(r), "=<>≤≥") {
		l.backup()
		return booleanOperatorState
	}
	if r == rune(EOF) {
		return eofState
	}

	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		l.backup()
		return elementState
	}

	panic(fmt.Errorf("unrecognised char %v, at pos: %d", string(r), l.pos))
}

func booleanOperatorState(l *Lexer) stFunc {
	l.addItem(BOOL_OPERATOR, string(l.next()))
	return startState
}

func conditionStartState(l *Lexer) stFunc {
	l.addItem(CONDITION_START, "[")
	return startState
}

func conditionEndState(l *Lexer) stFunc {
	l.addItem(CONDITION_END, "]")
	return startState
}

func delimeterState(l *Lexer) stFunc {

	r := l.next()

	if r == '/' {
		l.addItem(DELIMETER, "//")
		return startState
	} else {
		l.backup()
		l.addItem(DELIMETER, "/")
		return startState
	}
}

func attributeState(l *Lexer) stFunc {
	val := "@"
	for {
		r := l.next()
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			l.backup()
			break
		}
		val += string(r)
	}
	l.addItem(ATTRIBUTE, val)
	return startState
}

func elementState(l *Lexer) stFunc {

	var val string

	for {
		r := l.next()
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			l.backup()
			break
		}
		val += string(r)
	}
	l.addItem(ELEMENT, val)
	return startState
}

func eofState(l *Lexer) stFunc {
	l.addItem(EOF, "")
	return nil
}

// Helperfunctions

func (l *Lexer) next() rune {

	if l.pos >= len(l.input) {
		return rune(EOF)
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])

	//log.Printf("Read %v at pos %d", string(r), l.pos)

	l.width += w
	l.pos += w

	return r
}

func (l *Lexer) backup() {

	_, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width -= w
	l.pos -= w

	//log.Printf("Backed up. Pos is now %d", l.pos)
}

func (l *Lexer) addItem(it item_type, value string) {

	li := lexItem{typ: it, val: value, position: l.start}
	//log.Printf("Added item: %q", li)

	l.lexed_items = append(l.lexed_items, li)
	l.width = 0
	l.start = l.pos
}
