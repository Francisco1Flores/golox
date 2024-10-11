package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

var HadError bool = false
var tokens []scanner.Token

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 && len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename> or ./your_program.sh <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if !isCommandRight(command) {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	fileName := os.Args[2]
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		scan := scanner.NewScanner(fileContents)
		tokens = scan.Scan(fileContents)
		par := parser.NewParser(tokens)
		var expr *parser.Node

		switch command {
		case "run":
			stmt := par.ParseStmts()
			if errorHand.HadError {
				os.Exit(65)
			}
			inter := interpreter.NewStmtInterpreter(stmt)
			inter.ExecuteStmts()
		case "tokenize":
			scanner.PrintTokens(tokens)
		default:
			expr = par.ParseExpr()
			if command == "parse" {
				if !errorHand.HadError {
					parser.AstPrint(expr)
				}
			} else { // command to evaluate
				inter := interpreter.NewExprInterpreter(expr)
				result, err := inter.Interpret()
				if err != nil {
					os.Exit(70)
				}
				fmt.Println(result)
			}
		}
	} else {
		fmt.Println("EOF  null")
	}

	if errorHand.HadError {
		os.Exit(65)
	}
}

func isCommandRight(command string) bool {
	return command == "tokenize" || command == "parse" || command == "evaluate" || command == "run"
}

func thereIsCommand() bool {
	return len(os.Args) != 2
}
