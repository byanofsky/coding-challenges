// main.go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

const BUFFER_SIZE = 1024

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: my-cat <file>\n")
	}

	file := os.Args[1]

	f, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", file)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	b := make([]byte, BUFFER_SIZE)
	for {
		n, err := reader.Read(b)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		if n == 0 && errors.Is(err, io.EOF) {
			break
		}

		fmt.Fprintf(os.Stdin, "%s", b)
		if _, err := os.Stdout.Write(b); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
		}
	}
	// // Define command-line flags
	// name := flag.String("name", "World", "Name to greet")
	// flag.Parse()

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
