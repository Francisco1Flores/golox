package interpreter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type Interpreter struct {
	expr *parser.Node
}

type result struct {
	Value     string
	valueType scanner.TokenType
}

func NewInterpreter(expr *parser.Node) *Interpreter {
	return &Interpreter{expr: expr}
}

func (inter *Interpreter) Interpret() (string, error) {
	result, err := evaluate(inter.expr)
	return result.Value, err
}

func evaluate(expr *parser.Node) (result, error) {

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

func evaluateBinary(expr *parser.Node) (result, error) {
	left, err := evaluate(expr.Left)
	if err != nil {
		return result{}, err
	}

	right, err := evaluate(expr.Right)
	if err != nil {
		return result{}, err
	}

	// evaluate equality operators
	//if expr.Value.TokenType == scanner.EQUAL_EQUAL {
	//	if left.valueType != right.valueType {
	//		return result{"false", scanner.FALSE}, nil
	//	} else if left.Value != right.Value {
	//		return result{"false", scanner.FALSE}, nil
	//	}
	//	return result{"true", scanner.TRUE}, nil
	//} else if expr.Value.TokenType == scanner.BANG_EQUAL {
	//	if left.valueType != right.valueType {
	//		return result{"true", scanner.TRUE}, nil
	//	} else if left.Value != right.Value {
	//		return result{"true", scanner.TRUE}, nil
	//	}
	//	return result{"false", scanner.FALSE}, nil
	//}

	//if (expr.Value.TokenType == scanner.STAR ||
	//	expr.Value.TokenType == scanner.SLASH) &&
	//	(left.valueType != scanner.NUMBER || right.valueType != scanner.NUMBER) {
	//	errorHand.Error(expr.Value.Line, "Operands must be numbers.")
	//	return result{}, errors.New("operands must be numbers")
	//}

	//if left.valueType != scanner.NUMBER || right.valueType != scanner.NUMBER {
	//	return result{left.Value + right.Value, scanner.STRING}, nil
	//}

	var nLeft float64
	var nRight float64

	if areNumbers(left, right) {
		nLeft, err = strconv.ParseFloat(left.Value, 64)
		if err != nil {
			panic(err)
		}
		nRight, err = strconv.ParseFloat(right.Value, 64)
		if err != nil {
			panic(err)
		}
	}

	var res float64

	switch expr.Value.TokenType {
	case scanner.EQUAL_EQUAL:
		if left.valueType != right.valueType {
			return result{"false", scanner.FALSE}, nil
		} else if left.Value != right.Value {
			return result{"false", scanner.FALSE}, nil
		}
		return result{"true", scanner.TRUE}, nil
	case scanner.BANG_EQUAL:
		if left.valueType != right.valueType {
			return result{"true", scanner.TRUE}, nil
		} else if left.Value != right.Value {
			return result{"true", scanner.TRUE}, nil
		}
		return result{"false", scanner.FALSE}, nil
	case scanner.MINUS:
		if !areNumbers(left, right) {
			errorHand.Error(expr.Value.Line, "Operands must be numbers.")
			return result{}, errors.New("operands must be numbers")
		}
		res = nLeft - nRight
		return result{formatResultNum(res), scanner.NUMBER}, nil
	case scanner.STAR:
		if !areNumbers(left, right) {
			errorHand.Error(expr.Value.Line, "Operands must be numbers.")
			return result{}, errors.New("operands must be numbers")
		}
		res = nLeft * nRight
		return result{formatResultNum(res), scanner.NUMBER}, nil
	case scanner.SLASH:
		if !areNumbers(left, right) {
			errorHand.Error(expr.Value.Line, "Operands must be numbers.")
			return result{}, errors.New("operands must be numbers")
		}
		if nRight == 0 {
			return result{}, errors.New("division by cero")
		}
		res = nLeft / nRight
		return result{formatResultNum(res), scanner.NUMBER}, nil
	case scanner.PLUS:
		if areStrings(left, right) {
			return result{left.Value + right.Value, scanner.STRING}, nil
		} else if areNumbers(left, right) {
			res = nLeft + nRight
			return result{formatResultNum(res), scanner.NUMBER}, nil
		}
		errorHand.Error(expr.Value.Line, "Operands must be two numbers or two strings.")
		return result{}, errors.New("operands must be two numbers or two strings")
	case scanner.LESS:
		if nLeft < nRight {
			return result{"true", scanner.TRUE}, nil
		} else {
			return result{"false", scanner.FALSE}, nil
		}
	case scanner.LESS_EQUAL:
		if nLeft <= nRight {
			return result{"true", scanner.TRUE}, nil
		} else {
			return result{"false", scanner.FALSE}, nil
		}
	case scanner.GREATER:
		if nLeft > nRight {
			return result{"true", scanner.TRUE}, nil
		} else {
			return result{"false", scanner.FALSE}, nil
		}
	case scanner.GREATER_EQUAL:
		if nLeft >= nRight {
			return result{"true", scanner.TRUE}, nil
		} else {
			return result{"false", scanner.FALSE}, nil
		}
	}
	return result{}, errors.New("error in binary evaluation")
}

func evaluateUnary(expr *parser.Node) (result, error) {
	res, err := evaluate(expr.Right)

	if err != nil {
		return result{}, err
	}

	if expr.Value.Lexeme == "-" {
		if res.valueType != scanner.NUMBER {
			errorHand.Error(expr.Value.Line, "Operand must be a number.")
			return result{}, errors.New("operand must be a number")
		}

		if res.Value[0] == '-' {
			return result{res.Value[1:], scanner.NUMBER}, nil
		}
		return result{"-" + res.Value, scanner.NUMBER}, nil
	}

	if isTruthy(res.Value) {
		return result{"false", scanner.FALSE}, nil
	}
	return result{"true", scanner.TRUE}, nil
}

func evaluateLiteral(expr *parser.Node) (result, error) {
	switch expr.Value.TokenType {
	case scanner.TRUE:
		return result{"true", scanner.TRUE}, nil
	case scanner.FALSE:
		return result{"false", scanner.FALSE}, nil
	case scanner.NIL:
		return result{"nil", scanner.NIL}, nil
	case scanner.STRING:
		return result{expr.Value.Literal, scanner.STRING}, nil
	case scanner.NUMBER:
		return result{evaluateNumber(expr.Value.Literal), scanner.NUMBER}, nil
	}
	return result{}, errors.New("should not happend")
}

func evaluateGrouping(expr *parser.Node) (result, error) {
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

func formatResultNum(number float64) string {
	truncRes := int64(number)
	if number > float64(truncRes) {
		return fmt.Sprintf("%g", number)
	}
	return fmt.Sprintf("%.0f", number)
}

func isTruthy(value string) bool {
	if value == "nil" || value == "false" {
		return false
	}
	return true
}

func endsWith(text, patrn string) bool {
	pl := len(patrn)
	endtxt := len(text)
	return text[endtxt-pl:] == patrn
}

func areNumbers(left, right result) bool {
	return left.valueType == scanner.NUMBER && right.valueType == scanner.NUMBER
}

func areStrings(left, right result) bool {
	return left.valueType == scanner.STRING && right.valueType == scanner.STRING
}
