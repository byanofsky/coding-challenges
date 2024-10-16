package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Parse arguments
	flag.Parse()
	args := flag.Args()

	// TODO: Validate args input
	re := regexp.MustCompile(`^"(.*)"|(.*)$`)
	pattern := re.FindStringSubmatch(args[0])[0]
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
	// Track whether match is found
	foundMatch := false
	for scanner.Scan() {
		line := scanner.Text()

		// Check whether line matches
		if pattern != `""` { // Empty expression matches all
			matched, err := regexp.MatchString(pattern, line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Regexp error: %s\n", err)
			}
			if !matched {
				continue
			}
		}
		foundMatch = true

		// Add line breaks that were removed during scan.
		if firstLine {
			firstLine = false
		} else {
			w.WriteByte('\n')
		}

		w.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		// TODO: Extract to error checker function
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Exit code 1 when match not found.
	if !foundMatch {
		os.Exit(1)
	}
}
