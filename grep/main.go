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
	"strings"
)

func parsePattern(pattern string) string {
	p := strings.ReplaceAll(pattern, `\`, `\\`)
	// Attempt to unquote. If unquote has error, pattern likely doesn't have quotes.
	// So just return pattern as is.
	unquoted, err := strconv.Unquote(p)
	if err != nil {
		return pattern
	}

	return unquoted
}

func matchLine(pattern string, line string, caseInsensitive bool) bool {
	// Check whether line matches
	if len(pattern) == 0 { // Empty expression matches all
		return true
	}
	if caseInsensitive {
		pattern = strings.ToLower(pattern)
		line = strings.ToLower(line)
	}
	matched, err := regexp.MatchString(pattern, line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Regexp error: %s\n", err)
		os.Exit(1)
	}
	return matched
}

func grepFile(file string, pattern string, w *bufio.Writer, printFile bool, invertSearch bool, caseInsensitive bool) bool {
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
	// TODO: Swap to Reader.readline for longer lines
	// See: https://stackoverflow.com/a/21124415
	for scanner.Scan() {
		line := scanner.Text()
		isMatch := matchLine(pattern, line, caseInsensitive)
		if invertSearch {
			isMatch = !isMatch
		}
		if !isMatch {
			continue
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

func recurseGrep(w *bufio.Writer, invert bool, caseInsensitive bool) bool {
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: grep -r <pattern> <file>")
	}
	pattern := parsePattern(flag.Arg(0))
	root := flag.Arg(1)

	foundMatch := false

	// Walk the files from root. Skip dirs.
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if grepFile(path, pattern, w, true, invert, caseInsensitive) {
			foundMatch = true
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error recursing path: %s\n", err)
	}

	return foundMatch
}

func grep(w *bufio.Writer, invert bool, caseInsensitive bool) bool {
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: grep <pattern> <file>")
	}
	pattern := parsePattern(flag.Arg(0))
	file := flag.Arg(1)

	return grepFile(file, pattern, w, false, invert, caseInsensitive)
}

func main() {
	// Parse arguments
	recurse := flag.Bool("r", false, "Recurse a directory.")
	invert := flag.Bool("v", false, "Inverts search, excluding lines that match pattern.")
	caseInsensitive := flag.Bool("i", false, "Makes pattern case insensitive.")
	flag.Parse()

	// Create buffered writer to Stdout
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var foundMatch bool
	if *recurse {
		foundMatch = recurseGrep(w, *invert, *caseInsensitive)
	} else {
		foundMatch = grep(w, *invert, *caseInsensitive)
	}

	// Exit code 1 when match not found.
	if !foundMatch {
		os.Exit(1)
	}
}
