package interpreter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type Interpreter struct {
	expr *parser.Node
}

func NewInterpreter(expr *parser.Node) *Interpreter {
	return &Interpreter{expr: expr}
}

func (inter *Interpreter) Interpret() (string, error) {
	return evaluate(inter.expr)
}

func evaluate(expr *parser.Node) (string, error) {

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

func evaluateBinary(expr *parser.Node) (string, error) {
	left, err := evaluate(expr.Left)
	if err != nil {
		return "", err
	}
	right, err := evaluate(expr.Right)
	if err != nil {
		return "", err
	}

	if !isNumber(left) || !isNumber(right) {
		return left + right, nil
	}

	nLeft, err := strconv.ParseFloat(left, 64)
	if err != nil {
		panic(err)
	}
	nRight, err := strconv.ParseFloat(right, 64)
	if err != nil {
		panic(err)
	}

	var result float64

	switch expr.Value.Lexeme {
	case "-":
		result = nLeft - nRight
		//return fmt.Sprintf("%f", nLeft-nRight), nil
	case "*":
		result = nLeft * nRight
		//return fmt.Sprintf("%f", nLeft*nRight), nil
	case "/":
		if nRight == 0 {
			return "", errors.New("division by cero")
		}
		result = nLeft / nRight
	case "+":
		result = nLeft + nRight
	}
	truncRes := int64(result)
	if result > float64(truncRes) {
		return fmt.Sprintf("%g", result), nil
	}
	return fmt.Sprintf("%.0f", result), nil
}

func evaluateLiteral(expr *parser.Node) (string, error) {
	switch expr.Value.TokenType {
	case scanner.TRUE:
		return "true", nil
	case scanner.FALSE:
		return "false", nil
	case scanner.NIL:
		return "nil", nil
	case scanner.STRING:
		return expr.Value.Literal, nil
	case scanner.NUMBER:
		return evaluateNumber(expr.Value.Literal), nil
	}
	return "", nil
}

func evaluateUnary(expr *parser.Node) (string, error) {
	result, err := evaluate(expr.Right)

	if err != nil {
		return "", err
	}

	if expr.Value.Lexeme == "-" {
		if isNumber(result) {
			if result[0] == '-' {
				return result[1:], nil
			}
		}

		return "-" + result, nil
	}

	if isTruthy(result) {
		return "false", nil
	} else {
		return "true", nil
	}
}

func evaluateGrouping(expr *parser.Node) (string, error) {
	return evaluate(expr.Left)
}

func evaluateNumber(number string) string {
	if strings.Contains(number, ".") {
		if endsWith(number, ".0") {
			return number[:len(number)-2]
		}
	}
	return number
}

func isTruthy(value string) bool {
	if value == "nil" || value == "false" {
		return false
	}
	return true
}

func isNumber(number string) bool {
	for _, n := range number {
		if !isDigit(n) && n != '.' && n != '-' {
			return false
		}
	}
	return true
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func endsWith(text, patrn string) bool {
	pl := len(patrn)
	endtxt := len(text)
	return text[endtxt-pl:] == patrn
}
