package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func recurseGrep() bool {
	args := flag.Args()
	// TODO: Validate args input
	re := regexp.MustCompile(`^"(.*)"|(.*)$`)
	pattern := re.FindStringSubmatch(args[0])[0]
	file := args[1]
	fmt.Printf("file: %s pattern: %s", file, pattern)

	// TODO: Return foundMatch
	return true
}

func grep() bool {
	args := flag.Args()

	// TODO: Validate args input
	// TODO: use strcnv package to extract from within quotes
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

		// Write output
		out := fmt.Sprintf("%s\n", line)
		n, err := w.WriteString(out)
		if n != len(out) || err != nil {
			fmt.Fprintf(os.Stderr, "Error writing out: %s\n", err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		// TODO: Extract to error checker function
		fmt.Fprintf(os.Stderr, "Error printing: %s\n", err)
		os.Exit(1)
	}

	return foundMatch
}

func main() {
	// Parse arguments
	recurse := flag.Bool("r", false, "Recurse a directory. Usage: my-grp -r <pattern> <directory>")
	flag.Parse()

	var foundMatch bool
	if *recurse {
		foundMatch = recurseGrep()
	} else {
		foundMatch = grep()
	}

	// Exit code 1 when match not found.
	if !foundMatch {
		os.Exit(1)
	}
}
