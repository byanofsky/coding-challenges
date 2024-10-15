// main.go
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func openFile(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file does not exist: %s", file)
		}
		return nil, err
	}
	return f, err
}

func main() {
	numberLines := flag.Bool("n", false, "Number all output lines")
	numberNonBlankLines := flag.Bool("b", false, "Number non-blank output lines")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: my-cat [-n | -b] [file ...]")
		os.Exit(1)
	}

	if *numberLines && *numberNonBlankLines {
		fmt.Fprintln(os.Stderr, "Cannot provide both -n and -b flags")
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	lineNumber := 1
	for _, file := range flag.Args() {
		if err := processFile(file, w, numberLines, numberNonBlankLines, &lineNumber); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", file, err)
			os.Exit(1)
		}
	}
}

func processFile(file string, w *bufio.Writer, numberLines, numberNonBlankLines *bool, lineNumber *int) error {
	var r io.Reader

	if file == "-" {
		r = os.Stdin
	} else {
		f, err := openFile(file)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	scanner := bufio.NewScanner(r)
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
		} else {
			fmt.Fprintf(w, "\n")
		}
		text := scanner.Text()
		if *numberLines || (*numberNonBlankLines && len(text) > 0) {
			fmt.Fprintf(w, "%d  ", *lineNumber)
			*lineNumber++
		}
		fmt.Fprint(w, text)
	}
	return scanner.Err()
}
