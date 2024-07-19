package main

import "fmt"

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
	tokens = append(tokens, Token{line, "", "EOF"})
}

func scanTokens() {
	var c byte = advance()

	switch c {
	case '{':
		addToken(line, "{", "LEFT_BRACES")
	case '}':
		addToken(line, "}", "RIGHT_BRACES")
	case '(':
		addToken(line, "(", "LEFT_PAREN")
	case ')':
		addToken(line, ")", "RIGHT_PAREN")
	default:
	}
}

func PrintTokens() {
	for _, token := range tokens {
		var output string = fmt.Sprintf("%s %s %s", token.tokenType, token.value, "null")
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
