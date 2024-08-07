package main

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	// one character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	QUESTION_MARK
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	COLON
	SLASH
	STAR
	// one or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	// literals
	IDENTIFIER
	STRING
	NUMBER
	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

func (tokenType TokenType) String() string {
	return [42]string{
		"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE", "QUESTION_MARK",
		"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "COLON", "SLASH", "STAR", "BANG",
		"BANG_EQUAL", "EQUAL", "EQUAL_EQUAL", "GREATER", "GREATER_EQUAL", "LESS",
		"LESS_EQUAL", "IDENTIFIER", "STRING", "NUMBER", "AND", "CLASS", "ELSE", "FALSE",
		"FUN", "FOR", "IF", "NIL", "OR", "PRINT", "RETURN", "SUPER", "THIS", "TRUE", "VAR",
		"WHILE", "EOF"}[tokenType]
}

type Token struct {
	line      int
	lexeme    string
	literal   string
	tokenType TokenType
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
	eofToken := Token{line: line, tokenType: EOF, lexeme: "", literal: "null"}
	Tokens = append(Tokens, eofToken)
}

func scanTokens() {
	var c byte = advance()

	switch c {
	case '{':
		addToken(LEFT_BRACE)
	case '}':
		addToken(RIGHT_BRACE)
	case '(':
		addToken(LEFT_PAREN)
	case ')':
		addToken(RIGHT_PAREN)
	case ',':
		addToken(COMMA)
	case '.':
		addToken(DOT)
	case '-':
		addToken(MINUS)
	case '+':
		addToken(PLUS)
	case ';':
		addToken(SEMICOLON)
	case '*':
		addToken(STAR)
	case '\n':
		line++
	case '=':
		if match('=') {
			addToken(EQUAL_EQUAL)
		} else {
			addToken(EQUAL)
		}
	case '!':
		if match('=') {
			addToken(BANG_EQUAL)
		} else {
			addToken(BANG)
		}
	case '<':
		if match('=') {
			addToken(LESS_EQUAL)
		} else {
			addToken(LESS)
		}
	case '>':
		if match('=') {
			addToken(GREATER_EQUAL)
		} else {
			addToken(GREATER)
		}
	case '/':
		if match('/') {
			for !isAtEnd() && peek() != '\n' {
				advance()
			}
		} else {
			addToken(SLASH)
		}
	case ' ':
	case '\t':
	case '\r':
		break
	case '"':
		scanString()
	default:
		if isDigit(c) {
			scanNumber()
		} else if isAlpha(c) {
			scanIdentifier()
		} else {
			Error(line, "Unexpected character: "+string(c))
		}
	}
}

func PrintTokens() {
	for _, token := range Tokens {
		output := fmt.Sprintf("%s %s %s",
			token.tokenType.String(),
			token.lexeme,
			token.literal)
		fmt.Println(output)
	}
}

func scanString() {
	for !isAtEnd() && peek() != '"' {
		if peek() == '\n' {
			line++
		}
		advance()
	}
	if isAtEnd() {
		Error(line, "Unterminated string.")
		return
	}
	advance()

	value := source[start+1 : current-1]
	addTokenWithLiteral(STRING, value)
}

func scanNumber() {
	for isDigit(peek()) {
		advance()
	}
	if peek() == '.' && isDigit(peekNext()) {
		advance()
		for isDigit(peek()) {
			advance()
		}
	}
	fNumber, _ := strconv.ParseFloat(source[start:current], 64)
	sNumber := strconv.FormatFloat(fNumber, 'f', -1, 64)

	addTokenWithLiteral(NUMBER, sNumber)
}

func scanIdentifier() {

}

func advance() byte {
	current++
	return source[current-1]
}

func addToken(tokenType TokenType) {
	addTokenWithLiteral(tokenType, "null")
}

func addTokenWithLiteral(tokenType TokenType, literal string) {
	lexeme := source[start:current]
	tok := Token{
		line:      line,
		lexeme:    lexeme,
		tokenType: tokenType,
		literal:   literal,
	}

	Tokens = append(Tokens, tok)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isAtEnd() bool {
	return current >= len(source)
}

func peek() byte {
	if isAtEnd() {
		return 0
	}
	return source[current]
}

func peekNext() byte {
	if current+1 >= len(source) {
		return 0
	}
	return source[current+1]
}

func match(c byte) bool {
	if isAtEnd() {
		return false
	}
	if c != source[current] {
		return false
	}
	current++
	return true
}
