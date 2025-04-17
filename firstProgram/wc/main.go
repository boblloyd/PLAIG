package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	file := flag.String("f", "", "File to read as input")
	flag.Parse()

	i, err := run(*lines, *bytes, *file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(i)
}

func run(lines bool, bytes bool, file string) (int, error) {
	input := os.Stdin
	if file != "" {
		fileInput, err := os.Open(file)
		if err != nil {
			return 0, err
		}
		input = fileInput
		defer fileInput.Close()
	}
	return count(input, lines, bytes), nil
}

func count(r io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)
	if !countLines && !countBytes {
		scanner.Split(bufio.ScanWords)
	} else if !countLines {
		scanner.Split(bufio.ScanBytes)
	}
	wc := 0

	for scanner.Scan() {
		wc++
	}

	return wc
}
