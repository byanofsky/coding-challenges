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
	"strings"
)

type grepOptions struct {
	recurse         bool
	invert          bool
	caseInsensitive bool
	pattern         string
	path            string
}

func parsePattern(pattern string) string {
	p := pattern

	// Remove quotes
	if strings.HasPrefix(p, `'`) && strings.HasSuffix(p, `'`) || strings.HasPrefix(p, `"`) && strings.HasSuffix(p, `"`) {
		return p[1 : len(p)-1]
	}

	return p
}

func compilePattern(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

func matchLine(re *regexp.Regexp, line string) bool {
	return re.MatchString(line)
}

func grepFile(file string, re *regexp.Regexp, w *bufio.Writer, printFile, invert bool) (bool, error) {
	// Open file
	f, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("file does not exist: %s", file)
		}
		return false, fmt.Errorf("error opening file %s: %w", file, err)
	}
	defer f.Close()

	// Scan each line of file and print to Stdout buffered writer
	scanner := bufio.NewScanner(f)
	// Track whether match is found
	foundMatch := false

	for scanner.Scan() {
		line := scanner.Text()
		isMatch := matchLine(re, line)
		if invert {
			isMatch = !isMatch
		}
		if !isMatch {
			continue
		}
		foundMatch = true

		// Write output
		if printFile {
			fmt.Fprintf(w, "%s:%s\n", file, line)
		} else {
			fmt.Fprintln(w, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return foundMatch, fmt.Errorf("error reading file %s: %w", file, err)
	}

	return foundMatch, nil
}

func recurseGrep(opts grepOptions, w *bufio.Writer) (bool, error) {
	// TODO: Do this once
	re, err := compilePattern(opts.pattern, opts.caseInsensitive)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	foundMatch := false

	// Walk the files from root. Skip dirs.
	err = filepath.WalkDir(opts.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		match, err := grepFile(path, re, w, true, opts.invert)
		if err != nil {
			return err
		}
		if match {
			foundMatch = true
		}
		return nil
	})

	if err != nil {
		return foundMatch, fmt.Errorf("error recursing path: %w", err)
	}

	return foundMatch, nil
}

func grep(opts grepOptions, w *bufio.Writer) (bool, error) {
	// TODO: Compile pattern
	re, err := compilePattern(opts.pattern, opts.caseInsensitive)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return grepFile(opts.path, re, w, false, opts.invert)
}

func run() error {
	// Parse arguments
	recurse := flag.Bool("r", false, "Recurse a directory.")
	invert := flag.Bool("v", false, "Inverts search, excluding lines that match pattern.")
	caseInsensitive := flag.Bool("i", false, "Makes pattern case insensitive.")
	flag.Parse()

	if flag.NArg() != 2 {
		return errors.New("usage: my-grep [-r] [-v] [-i] <pattern> <file|irectory>")
	}

	opts := grepOptions{
		recurse:         *recurse,
		invert:          *invert,
		caseInsensitive: *caseInsensitive,
		pattern:         parsePattern(flag.Arg(0)),
		path:            flag.Arg(1),
	}

	// Create buffered writer to Stdout
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var foundMatch bool
	var err error

	if opts.recurse {
		foundMatch, err = recurseGrep(opts, w)
	} else {
		foundMatch, err = grep(opts, w)
	}

	if err != nil {
		return err
	}

	// Exit code 1 when match not found.
	if !foundMatch {
		return errors.New("no matches found")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
