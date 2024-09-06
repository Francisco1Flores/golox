package interpreter

import (
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type Interpreter struct {
	expr *parser.Node
}

func NewInterpreter(expr *parser.Node) *Interpreter {
	return &Interpreter{expr: expr}
}

func (inter *Interpreter) Interpret() string {
	return evaluate(inter.expr)
}

func evaluate(expr *parser.Node) string {

	switch expr.ExprType {
	case parser.BINARY:
		return evaluateBinary(expr)
	case parser.GROUPING:
		return evaluateGrouping(expr)
	case parser.UNARY:
		return evaluateUnary(expr)
	default:
		return evaluateLiteral(expr)
	}
}

func evaluateBinary(expr *parser.Node) string {

	return ""
}

func evaluateLiteral(expr *parser.Node) string {
	switch expr.Value.TokenType {
	case scanner.TRUE:
		return "true"
	case scanner.FALSE:
		return "false"
	case scanner.NIL:
		return "nil"
	case scanner.STRING:
		return expr.Value.Literal
	case scanner.NUMBER:
		return expr.Value.Literal
	}
	return ""
}

func evaluateUnary(expr *parser.Node) string {

	return ""
}

func evaluateGrouping(expr *parser.Node) string {

	return ""
}
