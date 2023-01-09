package interpreter

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

type Interpreter struct {
	writer io.Writer
}

func NewInterpreter(writer io.Writer) *Interpreter {
	return &Interpreter{
		writer: writer,
	}
}

// ----------------------------------------------------------------------------
// RuntimerError support

type RuntimeError struct {
	msg   string
	token *token.Token
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Interpreter Runtime Error: %s\nToken: %v\n", e.msg, e.token)
}

// ----------------------------------------------------------------------------
// Interpreter support

// Lox follows Ruby’s simple rule: false and nil are falsey, and everything else
// is truthy.
func (i *Interpreter) isTruthy(expr interface{}) bool {
	if expr == nil {
		return false
	} else if value, ok := expr.(bool); ok {
		return value
	} else {
		return true
	}
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil {
		return b == nil
	}

	switch ta := a.(type) {
	case string:
		tb, ok := b.(string)
		if !ok {
			return false
		}
		return ta == tb
	case float64:
		tb, ok := b.(float64)
		if !ok {
			return false
		}
		// According to IEEE 754, NaN is not equal to itself. But Java's .Equals
		// method makes all NaNs equal. JLox does the same, so we'll do the
		// same, for consistency.
		if math.IsNaN(ta) && math.IsNaN(tb) {
			return true
		}
		return ta == tb
	case bool:
		tb, ok := b.(bool)
		if !ok {
			return false
		}
		return ta == tb
	}
	panic("Implementation error: Interpreter.isEqual encountered a type that is not a string, float, or bool.")
}

func (i *Interpreter) checkNumberOperand(operator *token.Token, operand interface{}) (num *float64, err error) {
	if num, ok := operand.(float64); ok {
		return &num, nil
	} else {
		err = &RuntimeError{token: operator, msg: "Operand must be a number"}
		return nil, err
	}
}

// This is another of those pieces of code like isTruthy() that crosses the
// membrane between the user’s view of Lox objects and their internal
// representation in Java
func (i *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}

	if numVal, ok := val.(float64); ok {
		return strconv.FormatFloat(numVal, 'f', -1, 64)
	}

	return fmt.Sprintf("%s", val)

}

// ----------------------------------------------------------------------------
// Interpreter visitor

func (i *Interpreter) evaluate(expr ast.Expr) (result interface{}, err error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteral(expr *ast.LiteralExpr) (result interface{}, err error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGrouping(expr *ast.GroupingExpr) (result interface{}, err error) {
	return expr.Expression.Accept(i)
}

func (i *Interpreter) VisitUnary(expr *ast.UnaryExpr) (result interface{}, err error) {
	right, err := expr.Right.Accept(i)
	if err != nil {
		return
	}
	switch expr.Operator.TokenType {
	case token.Minus:
		result, err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -(*result), nil
	case token.Bang:
		result = !i.isTruthy(right)
	}
	return
}

func (i *Interpreter) VisitBinary(expr *ast.BinaryExpr) (result interface{}, err error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return
	}
	right, err := expr.Right.Accept(i)
	if err != nil {
		return
	}

	switch expr.Operator.TokenType {
	case token.Minus:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) - (*rightNum)
		return
	case token.Slash:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}

		if *rightNum == 0 {
			if *leftNum == 0 {
				result = math.NaN()
				return
			} else {
				err = &RuntimeError{
					msg:   fmt.Sprintf("Cannot divide by zero."),
					token: expr.Operator,
				}
				return
			}
		}
		result = (*leftNum) / (*rightNum)
		return
	case token.Star:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) * (*rightNum)
		return
	case token.Plus:
		numLeft, okLeft := left.(float64)
		numRight, okRight := right.(float64)
		if okLeft && okRight {
			result = numLeft + numRight
			return
		}

		strLeft, okLeft := left.(string)
		strRight, okRight := right.(string)
		if okLeft && okRight {
			result = strLeft + strRight
			return
		}

		err = &RuntimeError{
			msg:   fmt.Sprintf("operands must be both numbers or strings, got %v(%T) and %v(%T)", left, left, right, right),
			token: expr.Operator,
		}
		return
	case token.Greater:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) > (*rightNum)
		return
	case token.GreaterEqual:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) >= (*rightNum)
		return
	case token.Less:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) < (*rightNum)
		return
	case token.LessEqual:
		var leftNum *float64
		var rightNum *float64
		leftNum, err = i.checkNumberOperand(expr.Operator, left)
		if err != nil {
			return
		}
		rightNum, err = i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return
		}
		result = (*leftNum) <= (*rightNum)
		return
	case token.EqualEqual:
		result = i.isEqual(left, right)
		return
	case token.BangEqual:
		result = !i.isEqual(left, right)
		return

	}
	err = &RuntimeError{
		msg:   fmt.Sprintf("Unreachable code for operator."),
		token: expr.Operator,
	}
	return
}

// ----------------------------------------------------------------------------
// Interpreter API

func (i *Interpreter) Interpret(expr ast.Expr) error {
	value, err := expr.Accept(i)
	if err != nil {
		return &RuntimeError{
			msg: fmt.Sprintf("Interpreter failed exception: %v\n", err),
		}
	}
	fmt.Fprintln(i.writer, i.stringify(value))
	return nil
}
