package interpreter

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/internal/errorHand"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

type exprInterpreter struct {
	expr *parser.Node
}

type stmtInterpreter struct {
	stmts       []parser.Statement
	Environment *environment
}

/******************************************************************************/
type environment struct {
	values    map[string]result
	enclosing *environment
}

func newEnvironment() *environment {
	return &environment{
		values:    make(map[string]result),
		enclosing: nil,
	}
}

func (e *environment) define(name string, value result) {
	e.values[name] = value
}

func (e *environment) get(name string) result {
	value, ok := e.values[name]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	errorHand.Error(0, "Undefined variable '"+name+"'.")
	return result{}
}

func (e *environment) assign(name scanner.Token, value result) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
	errorHand.Error(0, "Undefined variable '"+name.Lexeme+"'.")
}

/******************************************************************************/

type result struct {
	Value     string
	valueType scanner.TokenType
}

func NewExprInterpreter(expr *parser.Node) *exprInterpreter {
	return &exprInterpreter{expr: expr}
}

func NewStmtInterpreter(stmts []parser.Statement) stmtInterpreter {
	return stmtInterpreter{
		stmts:       stmts,
		Environment: newEnvironment(),
	}
}

func (inter *exprInterpreter) Interpret() (string, error) {
	result, err := evaluate(inter.expr)
	return result.Value, err
}

func (s *stmtInterpreter) ExecuteStmts() {
	for _, stmt := range s.stmts {
		switch stmt.StmtType() {
		case parser.PRINT:
			stmt.Execute(func() {
				s.executePrintStmt(stmt)
			})
		case parser.EXPR:
			stmt.Execute(func() {
				s.executeExprStmt(stmt)
			})
		}
	}
}

func (s *stmtInterpreter) executePrintStmt(stmt parser.Statement) {
	pStmt, _ := stmt.(parser.PrintStmt)

	result, err := evaluate(pStmt.Expr)
	if err != nil {
		os.Exit(70)
	}
	fmt.Println(result.Value)
}

func (s *stmtInterpreter) executeExprStmt(stmt parser.Statement) {
	eStmt, _ := stmt.(parser.ExprStmt)
	_, err := evaluate(eStmt.Expr)
	if err != nil {
		os.Exit(70)
	}
}

func (s *stmtInterpreter) executeVarStmt(stmt parser.Statement) {
	vstmt, _ := stmt.(parser.VarDeclStmt)
	result := result{}
	if vstmt.Initializer != nil {
		result, _ = evaluate(vstmt.Initializer)
	}
	s.Environment.define(vstmt.Name.Lexeme, result)
}

func evaluate(expr *parser.Node) (result, error) {
	switch expr.ExprType {
	case parser.BINARY:
		return evaluateBinary(expr)
	case parser.GROUPING:
		return evaluateGrouping(expr)
	case parser.UNARY:
		return evaluateUnary(expr)
	case parser.VARIABLE:
		return evaluateVariable(expr)
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
		return booleanResult(isEqual(left, right)), nil
	case scanner.BANG_EQUAL:
		return booleanResult(!isEqual(left, right)), nil
	case scanner.MINUS:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		res = nLeft - nRight
		return result{formatResultNum(res), scanner.NUMBER}, nil
	case scanner.STAR:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		res = nLeft * nRight
		return result{formatResultNum(res), scanner.NUMBER}, nil
	case scanner.SLASH:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		if nRight == 0 {
			errorHand.Error(expr.Value.Line, "Division by zero.")
			return result{}, errors.New("division by zero")
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
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		return booleanResult(nLeft < nRight), nil
	case scanner.LESS_EQUAL:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		return booleanResult(nLeft <= nRight), nil
	case scanner.GREATER:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		return booleanResult(nLeft > nRight), nil
	case scanner.GREATER_EQUAL:
		if !checkAreNumbers(left, right, expr.Value.Line) {
			return result{}, errors.New("operands must be numbers")
		}
		return booleanResult(nLeft >= nRight), nil
	}
	return result{}, errors.New("error in binary evaluation")
}

func evaluateUnary(expr *parser.Node) (result, error) {
	res, err := evaluate(expr.Right)

	if err != nil {
		return result{}, err
	}

	if expr.Value.TokenType == scanner.MINUS {
		if res.valueType != scanner.NUMBER {
			errorHand.Error(expr.Value.Line, "Operand must be a number.")
			return result{}, errors.New("operand must be a number")
		}

		if res.Value[0] == '-' {
			return result{res.Value[1:], scanner.NUMBER}, nil
		}
		return result{"-" + res.Value, scanner.NUMBER}, nil
	}

	return booleanResult(!isTruthy(res.Value)), nil
}

func evaluateVariable(expr *parser.Node) (result, error) {

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

func isEqual(left, right result) bool {
	return left.valueType == right.valueType && left.Value == right.Value
}

func checkAreNumbers(left, right result, line int) bool {
	if !areNumbers(left, right) {
		errorHand.Error(line, "Operands must be numbers.")
		return false
	}
	return true
}

func booleanResult(value bool) result {
	if value {
		return result{"true", scanner.TRUE}
	}
	return result{"false", scanner.FALSE}
}
