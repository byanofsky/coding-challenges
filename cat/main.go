// main.go
package main

import (
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
	// numberLines := flag.Bool("n", false, "Enables numbered lines")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: my-cat <file>\n")
		os.Exit(1)
	}

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

		if _, err := io.Copy(os.Stdout, r); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying to stdout: %v\n", file)
			os.Exit(1)
		}
	}

	// // Use the parsed flag
	// fmt.Printf("Hello, %s!\n", *name)

	// // Example of handling subcommands
	// if len(os.Args) > 1 {
	// 	switch os.Args[1] {
	// 	case "version":
	// 		fmt.Println("v0.1.0")
	// 	case "help":
	// 		printHelp()
	// 	default:
	// 		fmt.Printf("Unknown command: %s\n", os.Args[1])
	// 		printHelp()
	// 	}
	// }
}

// func printHelp() {
// 	fmt.Println("Usage:")
// 	fmt.Println("  cli-tool [flags] [command]")
// 	fmt.Println("\nFlags:")
// 	flag.PrintDefaults()
// 	fmt.Println("\nCommands:")
// 	fmt.Println("  version    Print the version number")
// 	fmt.Println("  help       Print this help message")
// }
