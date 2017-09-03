package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <--key1=value1> <--key2=value2> <--flag1> ... template.file\n", os.Args[0])
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
		if len(i) >= 2 && i[0:2] == "--" {
			s := strings.Split(i[2:], "=")
			if len(s) > 1 {
				data[s[0]] = s[1]
			} else {
				data[i[2:]] = "true"
			}
		} else if templateFile == "" {
			templateFile = i
		} else {
			fmt.Fprintln(os.Stderr, "Error: more than one template file defined")
			usage()
		}
	}

	var templateContent string

	if templateFile == "" {
		fmt.Fprintln(os.Stderr, "Error: no template file defined")
		usage()
	} else if templateFile == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: can't read from STDIN")
			os.Exit(1)
		}
		templateContent = string(data)
	} else {
		// check template file
		if _, err := os.Stat(templateFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: can't find template file '%s'\n", templateFile)
			os.Exit(1)
		}
		// read from template
		data, err := ioutil.ReadFile(templateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: can't read from template file '%s'\n", templateFile)
			os.Exit(1)
		}
		templateContent = string(data)
	}

	// parse template file
	t, err := template.New("template").Parse(templateContent)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: can't parse template file")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	t.Option("missingkey=error")

	// render template file with environment variables
	var output bytes.Buffer
	err = t.Execute(&output, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: can't render template file")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// write result to STDOUT
	fmt.Print(output.String())
}
