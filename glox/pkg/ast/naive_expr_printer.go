package ast

import (
	"fmt"
	"strings"
)

type AstPrint struct {
}

func (a *AstPrint) print(e *Expr) (result interface{}, err error) {
	return (*e).Accept(a)
}

func (a *AstPrint) VisitBinary(e *BinaryExpr) (result interface{}, err error) {
	return a.parenthesize(e.operator.Lexeme, e.left, e.right)
}

func (a *AstPrint) VisitGrouping(e *GroupingExpr) (result interface{}, err error) {
	return a.parenthesize("group", e.expression)
}

func (a *AstPrint) VisitLiteral(e *LiteralExpr) (result interface{}, err error) {
	if e.value == nil {
		return "nil", nil
	} else {
		return fmt.Sprintf("%v", e.value), nil
	}
}

func (a *AstPrint) VisitUnary(e *UnaryExpr) (result interface{}, err error) {
	return a.parenthesize(e.operator.Lexeme, e.right)
}

func (a *AstPrint) parenthesize(name string, exprs ...Expr) (result string, err error) {
	builder := strings.Builder{}

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		exprStr, err := expr.Accept(a)
		if err != nil {
			return "", fmt.Errorf("Failed to stringify expr: name: %s, err: %w", name, err)
		}
		builder.WriteString(exprStr.(string))
	}
	builder.WriteString(")")
	return builder.String(), nil
}
