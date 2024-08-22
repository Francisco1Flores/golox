package parser

import (
	"errors"
	"fmt"
	"log"
	"strconv"

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
	switch expr.value.TokenType {
	case scanner.NUMBER:
		return stringifyNumber(expr.value.Lexeme)
	default:
		return expr.value.Lexeme
	}
}

func stringifyNumber(number string) string {
	numf, _ := strconv.ParseFloat(number, 64)
	trunCnum := float64(int32(numf))
	if numf > trunCnum {
		return strconv.FormatFloat(numf, 'g', 'g', 64)
	}
	return fmt.Sprintf("%.1f", numf)
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
	if parser.match(scanner.STRING, scanner.NUMBER) {
		return newNode(parser.previous(), nil, nil), nil
	}
	if parser.match(scanner.NIL) {
		return newNode(parser.previous(), nil, nil), nil
	}
	return Node{}, errors.New("expect expression")
}

func (parser *Parser) match(tokenType ...scanner.TokenType) bool {
	for _, tokt := range tokenType {
		if parser.tokens[parser.current].TokenType == tokt {
			(*parser).current++
			return true
		}
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
