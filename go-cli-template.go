package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func usage() {
	fmt.Printf("Usage: %s <--key1=value1> <--key2=value2> <--flag1> ... template.file\n", os.Args[0])
	os.Exit(1)
}

func main() {
	// get all environment variables
	data := make(map[string]string)
	for _, i := range os.Environ() {
		s := strings.Split(i, "=")
		data[s[0]] = s[1]
	}

	// parse command line arguments
	var templateFile string
	for _, i := range os.Args[1:] {
		if i[0:2] == "--" {
			s := strings.Split(i[2:], "=")
			if len(s) > 1 {
				data[s[0]] = s[1]
			} else {
				data[i[2:]] = "true"
			}
		} else if templateFile == "" {
			templateFile = i
		} else {
			fmt.Println("Error: more than one template file defined")
			usage()
		}
	}
	if templateFile == "" {
		fmt.Println("Error: no template file defined")
		usage()
	}

	// check template file
	if _, err := os.Stat(templateFile); err != nil {
		fmt.Printf("Error: can't find template file '%s'\n", templateFile)
		os.Exit(1)
	}

	// parse template file
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Println("Error: can't parse template file")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	t.Option("missingkey=error")

	// render template file with environment variables
	var output bytes.Buffer
	err = t.Execute(&output, data)
	if err != nil {
		fmt.Println("Error: can't render template file")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// write result to STDOUT
	fmt.Print(output.String())
}
