package scanner

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
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
	Line      int
	Lexeme    string
	Literal   string
	TokenType TokenType
}

type Scanner struct {
	source  []byte
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(sourceString []byte) *Scanner {
	return &Scanner{
		source:  sourceString,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

var keyWords map[string]TokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (scan *Scanner) Scan(sourceInput []byte) []Token {
	for !scan.isAtEnd() {
		scan.start = scan.current
		scan.scanTokens()
	}
	eofToken := Token{
		Line:      scan.line,
		TokenType: EOF,
		Lexeme:    "",
		Literal:   "null",
	}
	scan.tokens = append(scan.tokens, eofToken)
	return scan.tokens
}

func (scan *Scanner) scanTokens() {
	var c byte = scan.advance()

	switch c {
	case '{':
		scan.addToken(LEFT_BRACE)
	case '}':
		scan.addToken(RIGHT_BRACE)
	case '(':
		scan.addToken(LEFT_PAREN)
	case ')':
		scan.addToken(RIGHT_PAREN)
	case ',':
		scan.addToken(COMMA)
	case '.':
		scan.addToken(DOT)
	case '-':
		scan.addToken(MINUS)
	case '+':
		scan.addToken(PLUS)
	case ';':
		scan.addToken(SEMICOLON)
	case '*':
		scan.addToken(STAR)
	case '\n':
		scan.line++
	case '=':
		if scan.match('=') {
			scan.addToken(EQUAL_EQUAL)
		} else {
			scan.addToken(EQUAL)
		}
	case '!':
		if scan.match('=') {
			scan.addToken(BANG_EQUAL)
		} else {
			scan.addToken(BANG)
		}
	case '<':
		if scan.match('=') {
			scan.addToken(LESS_EQUAL)
		} else {
			scan.addToken(LESS)
		}
	case '>':
		if scan.match('=') {
			scan.addToken(GREATER_EQUAL)
		} else {
			scan.addToken(GREATER)
		}
	case '/':
		if scan.match('/') {
			for !scan.isAtEnd() && scan.peek() != '\n' {
				scan.advance()
			}
		} else {
			scan.addToken(SLASH)
		}
	case ' ':
	case '\t':
	case '\r':
		break
	case '"':
		scan.scanString()
	default:
		if isDigit(c) {
			scan.scanNumber()
		} else if isAlpha(c) || c == '_' {
			scan.scanIdentifier()
		} else {
			errorHand.Error(scan.line, "Unexpected character: "+string(c))
		}
	}
}

func PrintTokens(tokens []Token) {
	for _, token := range tokens {
		output := fmt.Sprintf("%s %s %s",
			token.TokenType.String(),
			token.Lexeme,
			token.Literal)
		fmt.Println(output)
	}
}

func (scan *Scanner) scanString() {
	for !scan.isAtEnd() && scan.peek() != '"' {
		if scan.peek() == '\n' {
			(*scan).line++
		}
		scan.advance()
	}
	if scan.isAtEnd() {
		errorHand.Error(scan.line, "Unterminated string.")
		return
	}
	scan.advance()

	value := string(scan.source[scan.start+1 : scan.current-1])
	scan.addTokenWithLiteral(STRING, value)
}

func (scan *Scanner) scanNumber() {
	var sNumber string
	for isDigit(scan.peek()) {
		scan.advance()
	}
	if scan.peek() == '.' && isDigit(scan.peekNext()) {
		scan.advance()
		for isDigit(scan.peek()) {
			scan.advance()
		}

		sNumber = string(scan.source[scan.start:scan.current])

		i := 1
		for sNumber[len(sNumber)-i] == '0' {
			i++
		}
		sNumber = sNumber[:len(sNumber)-i]

		scan.addTokenWithLiteral(NUMBER, sNumber)
		return
	}
	sNumber = string(scan.source[scan.start:scan.current]) + ".0"
	scan.addTokenWithLiteral(NUMBER, sNumber)
}

func (scan *Scanner) scanIdentifier() {
	for isAlphaNumeric(scan.peek()) || scan.peek() == '_' {
		scan.advance()
	}

	value := string(scan.source[scan.start:scan.current])
	tokenType, ok := keyWords[value]

	if !ok {
		tokenType = IDENTIFIER
	}
	scan.addToken(tokenType)
}

func (scan *Scanner) advance() byte {
	(*scan).current++

	return scan.source[scan.current-1]
}

func (scan *Scanner) addToken(tokenType TokenType) {
	scan.addTokenWithLiteral(tokenType, "null")
}

func (scan *Scanner) addTokenWithLiteral(tokenType TokenType, literal string) {
	lexeme := string(scan.source[scan.start:scan.current])
	tok := Token{
		Line:      scan.line,
		Lexeme:    lexeme,
		TokenType: tokenType,
		Literal:   literal,
	}

	scan.tokens = append(scan.tokens, tok)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (scan *Scanner) isAtEnd() bool {
	return scan.current >= len(scan.source)
}

func (scan *Scanner) peek() byte {
	if scan.isAtEnd() {
		return 0
	}
	return scan.source[scan.current]
}

func (scan *Scanner) peekNext() byte {
	if scan.current+1 >= len(scan.source) {
		return 0
	}
	return scan.source[scan.current+1]
}

func (scan *Scanner) match(c byte) bool {
	if scan.isAtEnd() {
		return false
	}
	if c != scan.source[scan.current] {
		return false
	}
	scan.current++
	return true
}
