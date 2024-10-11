// main.go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define command-line flags
	name := flag.String("name", "World", "Name to greet")
	flag.Parse()

	// Use the parsed flag
	fmt.Printf("Hello, %s!\n", *name)

	// Example of handling subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println("v0.1.0")
		case "help":
			printHelp()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  cli-tool [flags] [command]")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nCommands:")
	fmt.Println("  version    Print the version number")
	fmt.Println("  help       Print this help message")
}
