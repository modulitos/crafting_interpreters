// Code generated by generate_ast. DO NOT EDIT.
// Eg of Go's AST: https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/ast/ast.go#1

package ast

import (
	"github.com/modulitos/glox/pkg/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) (result interface{}, err error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type ExprVisitor interface {
	VisitAssign(e *AssignExpr) (result interface{}, err error)
	VisitBinary(e *BinaryExpr) (result interface{}, err error)
	VisitGrouping(e *GroupingExpr) (result interface{}, err error)
	VisitLiteral(e *LiteralExpr) (result interface{}, err error)
	VisitUnary(e *UnaryExpr) (result interface{}, err error)
	VisitVariable(e *VariableExpr) (result interface{}, err error)
}

type AssignExpr struct {
	Name  *token.Token
	Value Expr
}

func (e *AssignExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitAssign(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitBinary(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitGrouping(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator *token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitUnary(e)
}

type VariableExpr struct {
	Name *token.Token
}

func (e *VariableExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitVariable(e)
}

type StmtVisitor interface {
	VisitExpression(e *ExpressionStmt) error
	VisitPrint(e *PrintStmt) error
	VisitVar(e *VarStmt) error
	VisitBlock(e *BlockStmt) error
	VisitIf(e *IfStmt) error
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpression(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrint(e)
}

type VarStmt struct {
	Name        *token.Token
	Initializer Expr
}

func (e *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVar(e)
}

type BlockStmt struct {
	Statements []Stmt
}

func (e *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlock(e)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitIf(e)
}
