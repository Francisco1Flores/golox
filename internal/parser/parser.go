package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type Parser struct {
	current int
	tokens  []scanner.Token
}

type ExprType int

const (
	LITERAL ExprType = iota
	UNARY
	BINARY
	GROUPING
)

func (e ExprType) toString() string {
	return []string{"LITERAL", "UNARY", "BINARY", "GROUPING"}[e]
}

type Node struct {
	value    scanner.Token
	exprType ExprType
	left     *Node
	right    *Node
}

func newNode(token scanner.Token, exprType ExprType, left, right *Node) *Node {
	return &Node{
		value:    token,
		exprType: exprType,
		left:     left,
		right:    right,
	}
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		current: 0,
		tokens:  tokens,
	}
}

func (parser *Parser) Parse() *Node {
	expr, err := parser.expression()

	if err != nil {
		errorHand.ParseError(parser.tokens[parser.current].Line, err.Error())
	}
	return expr
}

// ************************* AstPrinter section *************************

func AstPrint(expr *Node) {
	fmt.Println(stringify(expr))
}

func stringify(expr *Node) string {
	switch expr.exprType {
	case BINARY:
		return parenthesize(stringifyBinary(expr))
	case LITERAL:
		if expr.value.TokenType == scanner.NUMBER {
			return stringifyNumber(expr.value.Lexeme)
		} else if expr.value.TokenType == scanner.STRING {
			return expr.value.Lexeme[1 : len(expr.value.Lexeme)-1]
		}
		return expr.value.Lexeme
	case GROUPING:
		return stringifyGroup(expr)
	case UNARY:
		return stringifyUnary(expr)
	default:
		return expr.value.Lexeme
	}
}

func stringifyBinary(expr *Node) string {
	return expr.value.Lexeme + " " + stringify(expr.left) + " " + stringify(expr.right)
}

func stringifyNumber(number string) string {
	numf, _ := strconv.ParseFloat(number, 64)
	truncNum := float64(int32(numf))
	if numf > truncNum {
		return fmt.Sprintf("%g", numf)
	}
	return fmt.Sprintf("%.1f", numf)
}

func stringifyGroup(expr *Node) string {
	if expr.left != nil {
		return "(group " + stringify(expr.left) + ")"
	}
	return ""
}

func stringifyUnary(expr *Node) string {
	return parenthesize(expr.value.Lexeme + " " + stringify(expr.right))
}

func parenthesize(text string) string {
	return "(" + text + ")"
}

// **********************************************************************

func (parser *Parser) expression() (*Node, error) {
	expr, err := parser.comparisson()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (parser *Parser) comparisson() (*Node, error) {
	expr, err := parser.term()

	if err != nil {
		return nil, err
	}

	for parser.match(scanner.LESS, scanner.LESS_EQUAL, scanner.GREATER, scanner.GREATER_EQUAL) {
		operator := parser.previous()
		right, err := parser.term()

		if err != nil {
			return nil, err
		}

		expr = newNode(operator, BINARY, expr, right)
	}

	return expr, nil
}

func (parser *Parser) term() (*Node, error) {
	expr, err := parser.factor()

	if err != nil {
		return nil, err
	}

	for parser.match(scanner.PLUS, scanner.MINUS) {
		operator := parser.previous()
		right, err := parser.factor()

		if err != nil {
			return nil, err
		}

		expr = newNode(operator, BINARY, expr, right)
	}

	return expr, nil
}

func (parser *Parser) factor() (*Node, error) {
	expr, err := parser.unary()

	if err != nil {
		return nil, err
	}

	for parser.match(scanner.SLASH, scanner.STAR) {
		operator := parser.previous()
		right, err := parser.unary()

		if err != nil {
			return nil, err
		}

		expr = newNode(operator, BINARY, expr, right)
	}

	return expr, nil
}

func (parser *Parser) unary() (*Node, error) {
	if parser.match(scanner.BANG, scanner.MINUS) {
		operator := parser.previous()
		expr, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return newNode(operator, UNARY, nil, expr), nil
	}

	expr, err := parser.primary()

	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (parser *Parser) primary() (*Node, error) {
	if parser.match(scanner.TRUE) {
		return newNode(parser.previous(), LITERAL, nil, nil), nil
	}
	if parser.match(scanner.FALSE) {
		return newNode(parser.previous(), LITERAL, nil, nil), nil
	}
	if parser.match(scanner.STRING, scanner.NUMBER) {
		return newNode(parser.previous(), LITERAL, nil, nil), nil
	}
	if parser.match(scanner.NIL) {
		return newNode(parser.previous(), LITERAL, nil, nil), nil
	}
	if parser.match(scanner.LEFT_PAREN) {
		thisTok := parser.previous()
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		err = nil
		_, err = parser.consume(scanner.RIGHT_PAREN, "Expect ) after expression.")
		if err != nil {
			return nil, err
		}

		return newNode(thisTok, GROUPING, expr, nil), nil
	}

	return nil, errors.New("expect expression")
}

func (parser *Parser) match(tokenType ...scanner.TokenType) bool {
	for _, tokt := range tokenType {
		if parser.tokens[parser.current].TokenType == tokt {
			parser.current++
			return true
		}
	}
	return false
}

func (parser *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}
	return scanner.Token{}, errors.New(message)
}

func (parser *Parser) advance() scanner.Token {
	if !parser.isAtEnd() {
		parser.current++
		return parser.tokens[parser.current-1]
	}
	return parser.previous()
}

func (parser *Parser) peek() scanner.Token {
	return parser.tokens[parser.current]
}

func (parser Parser) previous() scanner.Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) check(tokenType scanner.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}
	return parser.peek().TokenType == tokenType
}

func (parser Parser) isAtEnd() bool {
	return parser.tokens[parser.current].TokenType == scanner.EOF
}
