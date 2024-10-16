package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse arguments
	flag.Parse()
	args := flag.Args()

	// TODO: Validate args input
	// exp := args[0]
	file := args[1]

	// Open file
	f, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", file)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Create buffered writer to Stdout
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// Scan each line of file and print to Stdout buffered writer
	scanner := bufio.NewScanner(f)
	firstLine := true
	for scanner.Scan() {
		// Add line breaks that were removed during scan.
		if firstLine {
			firstLine = false
		} else {
			w.WriteByte('\n')
		}

		text := scanner.Text()
		w.WriteString(text)
	}
	if err := scanner.Err(); err != nil {
		// TODO: Extract to error checker function
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}
}
