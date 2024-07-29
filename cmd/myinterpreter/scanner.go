package main

import (
	"fmt"
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

var Tokens []Token

func Scan(sourceInput []byte) {
	source = string(sourceInput)
	for !isAtEnd() {
		start = current
		scanTokens()
	}
	Tokens = append(Tokens, Token{line, "", "EOF"})
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
	case '\n':
		line++
	case '=':
		if match('=') {
			addToken(line, "==", "EQUAL_EQUAL")
		} else {
			addToken(line, "=", "EQUAL")
		}
		break
		//agregar la funcion de detectar si es un igual solo o un igual mas otro simbolo
	default:
		Error(line, "Unexpected character: "+string(c))
	}
}

func PrintTokens() {
	for _, token := range Tokens {
		output := fmt.Sprintf("%s %s %s", token.tokenType, token.value, "null")
		fmt.Println(output)
	}
}

func advance() byte {
	current++
	return source[current-1]
}

func addToken(line int, value, tokenType string) {
	Tokens = append(Tokens, Token{line, value, tokenType})
}

func isAtEnd() bool {
	return current == len(source)
}

func peek() byte {
	return source[current]
}

func match(c byte) bool {
	if isAtEnd() {
		return false
	}
	return c == advance()
}
