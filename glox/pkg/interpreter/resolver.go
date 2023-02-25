package interpreter

import (
	"fmt"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

type Resolver struct {
	interpreter *Interpreter
	scopes      scopes
}

func NewResolver(interpreter *Interpreter) Resolver {
	return Resolver{
		interpreter: interpreter,
		scopes:      make(scopes, 0),
	}

}

////////////////////////////////////////////////////////////////////////////////
// API
////////////////////////////////////////////////////////////////////////////////

func (r *Resolver) ResolveStmts(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		err := r.resolveStmt(stmt)
		// aggregate errors?
		if err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// scopes
////////////////////////////////////////////////////////////////////////////////

type scopes []map[string]bool

func (s *scopes) pop() (map[string]bool, error) {
	last := s.peek()
	(*s) = (*s)[:len(*s)-1]
	return last, nil
}

func (s *scopes) push(scope map[string]bool) {
	// TODO: is the deref necessary?
	*s = append(*s, scope)
}

func (s *scopes) peek() map[string]bool {
	if s.isEmpty() {
		panic("called peek() when scopes stack is empty.")
	}
	return (*s)[len(*s)-1]
}

func (s *scopes) isEmpty() bool {
	return len(*s) == 0
}

////////////////////////////////////////////////////////////////////////////////
// Resolver private methods
////////////////////////////////////////////////////////////////////////////////

func (r *Resolver) resolveStmt(stmt ast.Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) resolveFunction(f *ast.FunctionStmt) error {
	r.beginScope()
	defer r.endScope()
	for _, param := range f.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}
	err := r.ResolveStmts(f.Body)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) beginScope() {
	// r.scopes = append(r.scopes, make(map[string]bool))
	r.scopes.push(make(map[string]bool))
}
func (r *Resolver) endScope() error {
	_, err := r.scopes.pop()
	return err
}

// Declaration adds the variable to the innermost scope so that it shadows any
// outer one and so that we know the variable exists. We mark it as "not ready
// yet" by binding its name to false in the scope map. The value associated with
// a key in the scope map represents whether or not we have finished resolving
// that variable's initializer.
func (r *Resolver) declare(name *token.Token) error {
	if r.scopes.isEmpty() {
		return nil
	}
	scope := r.scopes.peek()
	if ok, _ := scope[name.Lexeme]; ok {
		return fmt.Errorf("variable %s at line %d already exists in the scope.", name.Lexeme, name.Line)
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = true
}

// We start at the innermost scope and work outwards, looking in each map for a
// matching name. If we find the variable, we resolve it, passing in the number
// of scopes between the current innermost scope and the scope where the
// variable was found. So, if the variable was found in the current scope, we
// pass in 0. If it’s in the immediately enclosing scope, 1
//
// If we walk through all of the block scopes and never find the variable, we
// leave it unresolved and assume it’s global.
func (r *Resolver) resolveLocal(e ast.Expr, name string) {
	// var i int;
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name]; ok {
			r.interpreter.resolve(e, len(r.scopes)-1-i)
			return

		}
	}
}

// ----------------------------------------------------------------------------
// Resolver visitor

func (r *Resolver) VisitBlock(stmt *ast.BlockStmt) (err error) {
	r.beginScope()
	// TODO: does an error in a defer statement propagate?
	defer r.endScope()
	return r.ResolveStmts(stmt.Statements)
}

func (r *Resolver) VisitVar(stmt *ast.VarStmt) error {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitVariable(e *ast.VariableExpr) (interface{}, error) {
	if !r.scopes.isEmpty() {
		initialized, ok := r.scopes.peek()[e.Name.Lexeme]
		if ok && !initialized {
			// If the variable exists in the current scope but its value is false, that
			// means we have declared it but not yet defined it. We report that error.
			err := fmt.Errorf("Can't read local variable in its own initializer.")
			return nil, err
		}
	}
	r.resolveLocal(e, e.Name.Lexeme)
	return nil, nil
}

func (r *Resolver) VisitAssign(e *ast.AssignExpr) (interface{}, error) {
	r.resolveExpr(e.Value)
	r.resolveLocal(e, e.Name.Lexeme)
	return nil, nil
}

// Unlike variables, we define the name eagerly, before resolving the function’s
// body. This lets a function recursively refer to itself inside its own body.
func (r *Resolver) VisitFunction(stmt *ast.FunctionStmt) error {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	r.define(stmt.Name)
	err = r.resolveFunction(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) VisitExpression(stmt *ast.ExpressionStmt) error {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIf(stmt *ast.IfStmt) error {
	// When we resolve an if statement, there is no control flow. We resolve the
	// condition and both branches. Where a dynamic execution steps only into
	// the branch that is run, a static analysis is conservative—it analyzes any
	// branch that could be run. Since either one could be reached at runtime,
	// we resolve both.
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitPrint(stmt *ast.PrintStmt) error {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturn(stmt *ast.ReturnStmt) error {
	if stmt.Value != nil {
		err := r.resolveExpr(stmt.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitWhile(stmt *ast.WhileStmt) error {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	return r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitBinary(expr *ast.BinaryExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Right)
	return nil, err
}

func (r *Resolver) VisitCall(expr *ast.CallExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}
	for _, param := range expr.Args {
		err := r.resolveExpr(param)
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}

func (r *Resolver) VisitGrouping(expr *ast.GroupingExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Expression)
	return nil, err
}

func (r *Resolver) VisitLiteral(expr *ast.LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogical(expr *ast.LogicalExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Right)
	return nil, err
}

func (r *Resolver) VisitUnary(expr *ast.UnaryExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Right)
	return nil, err
}
