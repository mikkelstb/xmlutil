package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args[1:]
	l := NewLexer(arguments[0])
	l.run()

	for id, item := range l.lexed_items {
		fmt.Printf("ID: %3d, POS: %3d, TYPE: %-3v VAL: %v \n", id, item.position, item.typ, item.val)
	}
}

// The framework of the lexer and assosiated objects and functions

//Lexer, has information of the input string, and the lexed items
type Lexer struct {
	input string
	start int
	pos   int
	width int

	lexed_items []lexItem
}

// lexItem has a type and a value ex(ATTRIBUTE, "@id")
type lexItem struct {
	position int
	typ      item_type
	val      string
}

// item_type is an int
type item_type int

const (
	ELEMENT item_type = iota
	DELIMETER
	ATTRIBUTE
	VALUE
	ARIT_OPERATOR
	BOOL_OPERATOR
	CONDITION_START
	CONDITION_END
	EOF
	ERROR
)

// A stFunc function takes a lexer as input and returns a stFunc
type stFunc func(*Lexer) stFunc

//Constructor for lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:       input,
		lexed_items: make([]lexItem, 0),
	}
	return l
}

func (l *Lexer) run() []lexItem {
	for state := startState; state != nil; {
		state = state(l)
	}
	return l.lexed_items
}

func (li lexItem) String() string {
	switch li.typ {
	case EOF:
		return "EOF"
	}
	// if len(li.val) > 10 {
	// 	return fmt.Sprintf("%. 10q...", li.val)
	// }
	return fmt.Sprintf("Pos: %d %v %v", li.position, li.typ, li.val)
}
