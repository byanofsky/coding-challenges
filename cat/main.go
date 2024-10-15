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

const BUFFER_SIZE = 1024

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
	// Define command-line flags
	numberLines := flag.Bool("n", false, "Enables numbered lines")
	numberLinesSkipBlank := flag.Bool("b", false, "Enables numbered lines")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: my-cat <file>\n")
		os.Exit(1)
	}

	if *numberLines && *numberLinesSkipBlank {
		fmt.Fprintf(os.Stderr, "Cannot provide both -n and -b flags\n")
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	lineNumber := 1
	for _, file := range args {
		var r io.Reader

		if file == "-" {
			r = os.Stdin
		} else {
			f, err := openFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			r = f
		}

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)
		firstLine := true
		for scanner.Scan() {
			if firstLine {
				firstLine = false
			} else {
				fmt.Fprintf(w, "\n")
			}
			text := scanner.Text()
			if *numberLines {
				fmt.Fprintf(w, "%d  ", lineNumber)
				lineNumber++
			} else if *numberLinesSkipBlank {
				if len(text) > 0 {
					fmt.Fprintf(w, "%d  ", lineNumber)
					lineNumber++
				}
			}
			fmt.Fprintf(w, "%s", text)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid input: %s", err)
			os.Exit(1)
		}
	}

}
