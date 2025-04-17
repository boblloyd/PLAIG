package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var list []string
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	flag.Func("f", "Comma separated list of files to read as input", func(value string) error {
		for _, str := range strings.Split(value, ",") {
			if len(str) > 0 {
				list = append(list, str)
			}
		}
		return nil
	})
	flag.Parse()

	i, err := run(*lines, *bytes, list)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(i)
}

func run(lines bool, bytes bool, files []string) (int, error) {
	input := os.Stdin
	if len(files) > 0 {
		counted := 0
		for i := 0; i < len(files); i++ {
			fileInput, err := os.Open(files[i])
			if err != nil {
				return 0, err
			}
			input = fileInput
			defer fileInput.Close()
			counted += count(input, lines, bytes)
		}
		return counted, nil
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
