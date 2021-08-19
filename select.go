package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
)

type xpathType int

const (
	ELEMENT xpathType = iota
	ATTRIBUTE
)

func selectFromXpath(xpath string, xml_file string) {

	//Open file location, ang get an io.Reader compatible object
	file_reader, err := os.Open(xml_file)
	if err != nil {
		panic(err)
	}

	//Get an XML node
	xml_node, err := xmlquery.Parse(file_reader)
	if err != nil {
		panic(err)
	}

	list, err := xmlquery.QueryAll(xml_node, xpath)
	if err != nil {
		panic(err)
	}

	switch analyseXpath(xpath) {
	case ELEMENT:
		for id, result := range list {
			fmt.Printf("ID: %v: %v \n", id, strings.Trim(result.InnerText(), " \n"))
		}
	case ATTRIBUTE:
		for id, result := range list {
			fmt.Printf("ID: %v: %v \n", id, result.Attr)
		}
	}
}

func analyseXpath(xpath string) xpathType {
	attribute_pattern := regexp.MustCompile(`@\w+?$`)

	if attribute_pattern.MatchString(xpath) {
		return ATTRIBUTE
	}
	return ELEMENT
}
