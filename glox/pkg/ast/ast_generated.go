// Code generated by generate_ast. DO NOT EDIT.
// Eg of Go's AST: https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/ast/ast.go#1

package ast

import (
	"github.com/modulitos/glox/pkg/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) (result interface{}, err error)
}

type ExprVisitor interface {
	VisitBinary(e *BinaryExpr) (result interface{}, err error)
	VisitGrouping(e *GroupingExpr) (result interface{}, err error)
	VisitLiteral(e *LiteralExpr) (result interface{}, err error)
	VisitUnary(e *UnaryExpr) (result interface{}, err error)
}

type BinaryExpr struct {
	left     Expr
	operator *token.Token
	right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitBinary(e)
}

type GroupingExpr struct {
	expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitGrouping(e)
}

type LiteralExpr struct {
	value interface{}
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	operator *token.Token
	right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitUnary(e)
}
