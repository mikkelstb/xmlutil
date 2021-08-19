package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func listElements(filename string, attribute bool) {

	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	decoder := xml.NewDecoder(reader)
	var stack []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Exit(1)
		}

		switch token := token.(type) {
		case xml.StartElement:
			stack = append(stack, token.Name.Local)
			fmt.Printf("Element: \t /%s\n", strings.Join(stack, "/"))
			if attribute {
				for _, attribute := range token.Attr {
					fmt.Printf("Attribute: \t /%s/@%s\n", strings.Join(stack, "/"), attribute.Name.Local)
				}
			}
		case xml.EndElement:
			stack = stack[0 : len(stack)-1]
		case xml.CharData:
			break
		default:
			break
		}
	}
}
