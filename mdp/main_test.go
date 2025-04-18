package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	inputFile    = "./testdata/test1.md"
	goldenFile   = "./testdata/test1.md.html"
	templateFile = "./testdata/test2.html.tmpl"
	goldenFile2  = "./testdata/test2.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result, err := parseContent(input, "", inputFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer
	if err := run(inputFile, "", &mockStdOut, true, false); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("result content does not match golden file")
	}

	os.Remove(resultFile)
}

func TestRunTemplate(t *testing.T) {
	var mockStdOut bytes.Buffer
	if err := run(inputFile, templateFile, &mockStdOut, true, false); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(goldenFile2)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("result content does not match golden file")
	}

	os.Remove(resultFile)
}
