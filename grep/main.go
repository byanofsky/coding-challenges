package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func parsePattern(pattern string) string {
	// Attempt to unquote. If unquote has error, pattern likely doesn't have quotes.
	// So just return pattern as is.
	unquoted, err := strconv.Unquote(pattern)
	if err != nil {
		return pattern
	}

	return unquoted
}

func grepFile(file string, pattern string, w *bufio.Writer, printFile bool) bool {
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

	// Scan each line of file and print to Stdout buffered writer
	scanner := bufio.NewScanner(f)
	// Track whether match is found
	foundMatch := false
	for scanner.Scan() {
		line := scanner.Text()

		// Check whether line matches
		if len(pattern) > 0 { // Empty expression matches all
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
		out := line
		if printFile {
			out = fmt.Sprintf("%s:%s", file, out)
		}
		out = fmt.Sprintf("%s\n", out)

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

func recurseGrep(w *bufio.Writer) bool {
	args := flag.Args()
	// TODO: Validate args input
	// TODO: Use flag.Arg(0)
	pattern := parsePattern(args[0])

	root := args[1]

	foundMatch := false

	// Walk the files from root. Skip dirs.
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if grepFile(path, pattern, w, true) {
			foundMatch = true
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error recursing path: %s\n", err)
	}

	return foundMatch
}

func grep(w *bufio.Writer) bool {
	args := flag.Args()

	// TODO: Validate args input length
	// TODO: Use flag.Arg(0)
	pattern := parsePattern(args[0])
	file := args[1]

	return grepFile(file, pattern, w, false)
}

func main() {
	// Parse arguments
	recurse := flag.Bool("r", false, "Recurse a directory. Usage: my-grp -r <pattern> <directory>")
	flag.Parse()

	// Create buffered writer to Stdout
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var foundMatch bool
	if *recurse {
		foundMatch = recurseGrep(w)
	} else {
		foundMatch = grep(w)
	}

	// Exit code 1 when match not found.
	if !foundMatch {
		os.Exit(1)
	}
}
