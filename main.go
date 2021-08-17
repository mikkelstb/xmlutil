package main

import (
	"os"
)


func main()  {
	
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		panic("to few arguments!")
	}
	command := arguments[0]

	switch(command) {
	case "el":
		listElements(arguments[1], false)
		break
	case "ela":
		listElements(arguments[1], true)
		break
	case "sel":
		selectFromXpath(arguments[1], arguments[2])
	}

}