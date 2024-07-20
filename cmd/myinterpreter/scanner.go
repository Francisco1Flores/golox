package main

import (
	"fmt"
	"os"
)

type Token struct {
	line      int
	value     string
	tokenType string
}

var start int = 0
var current int = 0
var line int = 1

var source string

var tokens []Token

func Scan(sourceInput []byte) {
	source = string(sourceInput)
	for !isAtEnd() {
		start = current
		scanTokens()
	}
	if len(tokens) == 0 && len(source) != 0 {
		os.Exit(65)
	}
	tokens = append(tokens, Token{line, "", "EOF"})
}

func scanTokens() {
	var c byte = advance()

	switch c {
	case '{':
		addToken(line, "{", "LEFT_BRACE")
	case '}':
		addToken(line, "}", "RIGHT_BRACE")
	case '(':
		addToken(line, "(", "LEFT_PAREN")
	case ')':
		addToken(line, ")", "RIGHT_PAREN")
	case ',':
		addToken(line, ",", "COMMA")
	case '.':
		addToken(line, ".", "DOT")
	case '-':
		addToken(line, "-", "MINUS")
	case '+':
		addToken(line, "+", "PLUS")
	case ';':
		addToken(line, ";", "SEMICOLON")
	case '*':
		addToken(line, "*", "STAR")
	default:
		ReportError(line, "Unexpected character: "+string(c))
	}
}

func PrintTokens() {
	if len(tokens) == 0 && len(source) != 0 {
		os.Exit(65)
	}
	for _, token := range tokens {
		output := fmt.Sprintf("%s %s %s", token.tokenType, token.value, "null")
		fmt.Println(output)
	}
}

func advance() byte {
	current++
	return source[current-1]
}

func addToken(line int, value, tokenType string) {
	tokens = append(tokens, Token{line, value, tokenType})
}

func isAtEnd() bool {
	return current == len(source)
}
