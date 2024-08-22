package parser

import (
	"errors"
	"fmt"
	"log"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type Parser struct {
	current int
	tokens  []scanner.Token
}

type Node struct {
	value scanner.Token
	left  *Node
	right *Node
}

func newNode(token scanner.Token, left, right *Node) Node {
	return Node{
		value: token,
		left:  left,
		right: right,
	}
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		current: 0,
		tokens:  tokens,
	}
}

func (parser Parser) Parse() Node {
	expr, err := parser.expression()

	if err != nil {
		log.Fatal(err)
		errorHand.ParseError(parser.tokens[parser.current].Line, "Error")
	}
	return expr
}

func AstPrint(expr Node) {
	fmt.Println(stringify(expr))
}

func stringify(expr Node) string {
	return expr.value.Lexeme
}

func (parser Parser) expression() (Node, error) {
	expr, err := parser.literal()
	if err != nil {
		return Node{}, err
	}
	return expr, nil
}

func (parser Parser) literal() (Node, error) {
	if parser.match(scanner.TRUE) {
		return newNode(parser.previous(), nil, nil), nil
	}
	if parser.match(scanner.FALSE) {
		return newNode(parser.previous(), nil, nil), nil
	}
	if parser.match(scanner.NIL) {
		return newNode(parser.previous(), nil, nil), nil
	}
	return Node{}, errors.New("expect expression")
}

func (parser *Parser) match(tokenType scanner.TokenType) bool {
	if parser.tokens[parser.current].TokenType == tokenType {
		(*parser).current++
		return true
	}
	return false
}

func (parser Parser) peek() scanner.Token {
	return parser.tokens[parser.current]
}

func (parser Parser) previous() scanner.Token {
	return parser.tokens[parser.current-1]
}

func (parser Parser) isAtEnd() bool {
	return parser.tokens[parser.current].TokenType == scanner.EOF
}
