package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

type Property struct {
	propShortName   string
	propType        string
	propName        string
	propDescription string
}

func main() {
	csvFile, _ := os.Open("properties.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	out := os.Stdout

	var shortName string
	var typeName string

	var properties []Property
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		shortName = getShortName(line[1])
		typeName = getTypeName(line[2])

		io.WriteString(out, shortName+"\t\t\t\t"+typeName+"\t\t\t\t"+line[1]+"\t\t\t\t//"+line[4]+"\n")

		properties = append(properties, Property{
			propShortName:   line[0],
			propType:        line[1],
			propName:        line[2],
			propDescription: line[3],
		})
	}
}

func getShortName(name string) string {
	if len(name) > 9 && name[:9] == "PROP_CAR_" {
		return strings.ToLower(name[9:len(name)])
	}
	if len(name) > 5 && name[:5] == "PROP_" {
		return strings.ToLower(name[5:len(name)])
	}

	return strings.ToLower(name)
}

func getTypeName(name string) string {

	switch name {
	case "tokenize":
		return "token"
	case "exist":
		return "bool"
	default:
		return name
	}
}
