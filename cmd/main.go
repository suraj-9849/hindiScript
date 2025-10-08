package main

import (
	"fmt"
	"github.com/suraj-9849/hindiLang.git/config/functions"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a .hlang file to run")
			fmt.Println("Usage: ./hlang.exe run <filename.hlang>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "version", "-v", "--version":
		fmt.Println("HindiScript v1.0.0")
		fmt.Println("A programming language in Hindi")
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	ext := filepath.Ext(filename)
	if ext != ".hlang" {
		fmt.Printf("Warning: File '%s' does not have .hlang extension\n", filename)
	}

	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	err = functions.Run(string(code))
	if err != nil {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("HindiScript - A programming language in Hindi")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("./hlang.exe run <filename.hlang>    Run a .hlang file")
	fmt.Println("./hlang.exe version                 Show version information")
	fmt.Println("./hlang.exe help                    Show this help message")
	fmt.Println()
}