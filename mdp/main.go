package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="content-type" content="text/html; charset=utf-8">
        <title>{{ .Title }}</title>
    </head>
    <body>
<h1>{{ .Filename }}</h1>

{{ .Body }}
    </body>
</html>
`
)

type content struct {
	Title    string
	Body     template.HTML
	Filename string
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview.  Do not use with -i")
	skipPreview := flag.Bool("s", false, "Skip auto-previewing the file after conversion")
	tFname := flag.String("t", "", "Alternate template name")
	useStdin := flag.Bool("i", false, "Use standard in for input.  Do not use with -file")
	flag.Parse()

	if *filename == "" && !*useStdin {
		flag.Usage()
		os.Exit(1)
	}
	if *filename != "" && *useStdin {
		flag.Usage()
		os.Exit(1)
	}

	if os.Getenv("MDP_TEMPLATE_FILENAME") != "" {
		*tFname = os.Getenv("MDP_TEMPLATE_FILENAME")
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview, *useStdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, tFname string, out io.Writer, skipPreview bool, useStdin bool) error {
	var input []byte
	var err error
	if useStdin {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		input, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
	}
	htmlData, err := parseContent(input, tFname, filename)
	if err != nil {
		return err
	}

	temp, err := ioutil.TempFile("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}

	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, tFname string, filename string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body),
		Filename: filename,
	}

	var buffer bytes.Buffer

	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFName string, data []byte) error {
	return ioutil.WriteFile(outFName, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	err = exec.Command(cPath, cParams...).Run()

	time.Sleep(2 * time.Second)
	return err
}
