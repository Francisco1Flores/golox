package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

var HadError bool = false
var tokens []scanner.Token

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		scan := scanner.NewScanner(fileContents)
		tokens = scan.Scan(fileContents)
		par := parser.NewParser(tokens)
		var expr parser.Node

		switch command {
		case "tokenize":
			scanner.PrintTokens(tokens)
		default:
			expr = par.Parse()
			parser.AstPrint(expr)
		}
	} else {
		fmt.Println("EOF  null")
	}

	if errorHand.HadError {
		os.Exit(65)
	}
}
