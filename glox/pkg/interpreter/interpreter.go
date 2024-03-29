package interpreter

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

// ----------------------------------------------------------------------------
// Interpreter API

type Interpreter struct {
	writer      io.Writer
	environment *environment // should this be a pointer?
	globals     *environment
	locals      map[ast.Expr]int
}

func NewInterpreter(writer io.Writer) *Interpreter {
	globals := newGlobalEnvironment()
	return &Interpreter{
		writer: writer,
		// Pointer to the current env, which can change as we traverse blocks:
		environment: globals,
		// Pointer to the global env:
		globals: globals,
		locals:  make(map[ast.Expr]int),
	}
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) error {

	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return &RuntimeError{
				msg: fmt.Sprintf("Interpreter failed exception: %v\n", err),
			}
		}
	}
	return nil
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

func (i *Interpreter) resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookupVariable(name *token.Token, expr ast.Expr) (interface{}, error) {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		return i.globals.get(name)
	}

}

// ----------------------------------------------------------------------------
// Interpreter visitor

func (i *Interpreter) execute(stmt ast.Stmt) error {
	return stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr ast.Expr) (result interface{}, err error) {
	return expr.Accept(i)
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, env *environment) (err error) {
	previous := i.environment
	i.environment = env
	i.environment.parent = previous
	defer func() {
		i.environment.parent = nil
		i.environment = previous
	}()

	for _, stmt := range stmts {
		err = i.execute(stmt)
		if err != nil {
			return
		}
	}
	return
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
		// Many languages define + such that if either operand is a string, the
		// other is converted to a string and the results are then concatenated.
		switch tl := left.(type) {
		case float64:
			switch tr := right.(type) {
			case float64:
				result = tl + tr
				return
			case string:
				result = fmt.Sprintf("%s%s", i.stringify(left), right)
				return
			}
		case string:
			switch tr := right.(type) {
			case float64:
				result = fmt.Sprintf("%s%s", left, i.stringify(right))
				return
			case string:
				result = tl + tr
				return
			}
		}

		err = &RuntimeError{
			msg:   fmt.Sprintf("operands must be both numbers, both strings, or at least one number and a string. Got %v(%T) and %v(%T)", left, left, right, right),
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

func (i *Interpreter) VisitCall(expr *ast.CallExpr) (result interface{}, err error) {
	var callee interface{}
	callee, err = i.evaluate(expr.Callee)
	if err != nil {
		return
	}
	var args []interface{}
	for _, argExpr := range expr.Args {
		var arg interface{}
		arg, err = i.evaluate(argExpr)
		if err != nil {
			return
		}
		args = append(args, arg)
	}
	function, ok := callee.(Callable)
	if !ok {
		err = &RuntimeError{
			msg:   fmt.Sprintf("Can only call functions and classes. Callee is unexpected type: %T", callee),
			token: expr.Paren,
		}
		return
	}
	if len(args) != function.arity() {
		err = &RuntimeError{
			msg:   fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(args)),
			token: expr.Paren,
		}
		return
	}
	result, err = function.call(i, args)
	if err != nil {
		return
	}
	return
}

func (i *Interpreter) VisitExpression(stmt *ast.ExpressionStmt) error {
	// Appropriately enough, we discard the value returned by i.evaluate() by
	// placing that call inside a Golang expression statement.
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrint(stmt *ast.PrintStmt) error {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Fprintln(i.writer, i.stringify(value))
	return nil
}

func (i *Interpreter) VisitVar(stmt *ast.VarStmt) error {
	var value interface{}
	var err error
	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
	}
	// We'll keep it simple and say that Lox sets a variable to nil if it isn’t
	// explicitly initialized.
	// how do we know whether to define this in the env or the global?
	i.environment.define(stmt.Name.Lexeme, value)
	return err
}

func (i *Interpreter) VisitVariable(e *ast.VariableExpr) (interface{}, error) {
	return i.lookupVariable(e.Name, e)
}

func (i *Interpreter) VisitAssign(e *ast.AssignExpr) (interface{}, error) {
	result, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}
	if distance, ok := i.locals[e]; ok {
		err := i.environment.assignAt(distance, e.Name.Lexeme, result)
		if err != nil {
			return nil, err
		}
	} else {
		i.globals.values[e.Name.Lexeme] = result
	}
	return result, nil
}

func (i *Interpreter) VisitBlock(stmt *ast.BlockStmt) (err error) {
	err = i.executeBlock(stmt.Statements, newEnvironment(i.environment))
	return
}

func (i *Interpreter) VisitIf(stmt *ast.IfStmt) (err error) {
	res, err := i.evaluate(stmt.Condition)
	if err != nil {
		return
	}
	if i.isTruthy(res) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return
}

func (i *Interpreter) VisitLogical(stmt *ast.LogicalExpr) (result interface{}, err error) {
	left, err := i.evaluate(stmt.Left)
	if err != nil {
		return
	}
	if stmt.Operator.TokenType == token.Or {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	result, err = i.evaluate(stmt.Right)
	if err != nil {
		return
	}
	return
}

func (i *Interpreter) VisitWhile(stmt *ast.WhileStmt) (err error) {
	for {
		var cond interface{}
		cond, err = i.evaluate(stmt.Condition)
		if err != nil {
			return
		}
		if i.isTruthy(cond) {
			err = i.execute(stmt.Body)
			if err != nil {
				return
			}
		} else {
			break
		}
	}
	return
}

func (i *Interpreter) VisitFunction(stmt *ast.FunctionStmt) (err error) {
	function := &loxFunction{
		declaration: stmt,
	}
	i.environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitReturn(stmt *ast.ReturnStmt) (err error) {
	var value interface{}
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return
		}
	}
	// unwind the stack:
	panic(&returnPayload{
		Value: value,
	})
}
