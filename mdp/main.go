package main

import {
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
}

const (
	header = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>Markdown Preview Tool</title>
	</head>
	<body>
`

	footer = `
	</body>
</html>
`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *filename == "" {
			flag.Usage()
			os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filname string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte)  []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SantizeBytes(output)

	var buffer bytes.buffer

	buffer.WriteString(header)
	buffer.WriteString(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFName string, data []byte) error {
	return ioutil.WriteFile(outFName, data, 0644)
}
