package main

import (
	"fmt"
	"mikkelstb/xmlutil/xpath"
	"os"
	"strings"
)

func main() {

	arguments := os.Args[1:]
	if len(arguments) < 1 {
		panic("to few arguments!")
	}
	command := arguments[0]

	switch command {
	case "el":
		listElements(arguments[1], false)
	case "ela":
		listElements(arguments[1], true)
	case "sel":
		selectFromXpath(arguments[1], arguments[2])
	case "lex":
		lexer := xpath.NewLexer(strings.NewReader(arguments[1]))

		for id, token := range lexer.LexAll() {
			fmt.Printf("ID: %v, POS: %v, TYPE: %v, VAL: %v \n", id, token.Position, token.Tokentype, token.Value)
		}
	}
}
