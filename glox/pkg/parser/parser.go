// Golang's parser:
// https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/parser/parser.go#524
package parser

import (
	"fmt"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

const maxFuncArgCounts = 255

type Parser struct {
	Tokens  []*token.Token
	current int
}

type ParserError struct {
	err error
}

func (e ParserError) Error() string {
	return fmt.Sprintf("ParserError: %v", e.err)
}

// ----------------------------------------------------------------------------
// Parsing support

func (p *Parser) peek() *token.Token {
	return p.Tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.Eof
}

func (p *Parser) previous() *token.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	} else {
		return tokenType == p.Tokens[p.current].TokenType
	}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType token.Type) (*token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	} else {
		actual := p.peek()
		return nil, fmt.Errorf("wanted token type: %s, got: %s, at line: %d", tokenType, actual, actual.Line)
	}
}

// Discards tokens until it think it has found a statement boundary.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.Semicolon {
			return
		}

		switch p.peek().TokenType {
		case token.Class, token.Fun, token.Var, token.For, token.If, token.While, token.Print, token.Return:
			return
		}

		p.advance()
	}
}

// ----------------------------------------------------------------------------
// Types

// declaration → funDecl | varDecl | statement;
func (p *Parser) declaration() (stmt ast.Stmt, err error) {
	defer func() {
		if err != nil {
			p.synchronize()
			return
		}
	}()

	if p.match(token.Fun) {
		return p.function("function")
	}
	if p.match(token.Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

// funDecl    → "fun" function ;
// function   → IDENTIFIER "(" parameters? ")" block ;
// parameters → IDENTIFIER ( "," IDENTIFIER )* ;
func (p *Parser) function(kind string) (stmt ast.Stmt, err error) {
	name, err := p.consume(token.Identifier)
	if err != nil {
		err = fmt.Errorf("Expect %s name: %w", kind, err)
		return
	}
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("Expect ( after %s name: %w", kind, err)
		return
	}

	var params []*token.Token
	if !p.check(token.RightParen) {
		for {
			if len(params) >= maxFuncArgCounts {
				err = fmt.Errorf("Can't have more than 255 parameters.")
				return
			}
			var param *token.Token
			param, err = p.consume(token.Identifier)
			if err != nil {
				err = fmt.Errorf("Expect parameter name: %w.", err)
				return
			}
			params = append(params, param)

			if !p.match(token.Comma) {
				break
			}
		}
	}
	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("Expect ) after %s parameters: %w", kind, err)
		return
	}
	_, err = p.consume(token.LeftBrace)
	if err != nil {
		err = fmt.Errorf("Expect { before %s body: %w", kind, err)
		return
	}
	body, err := p.block()
	if err != nil {
		return
	}
	return &ast.FunctionStmt{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil
}

func (p *Parser) varDeclaration() (stmt ast.Stmt, err error) {
	name, err := p.consume(token.Identifier)
	if err != nil {
		return
	}
	var initializer ast.Expr
	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return
		}
	}

	_, err = p.consume(token.Semicolon)
	if err != nil {
		return
	}
	return &ast.VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (stmt ast.Stmt, err error) {
	if p.match(token.For) {
		return p.forStatement()
	}
	if p.match(token.If) {
		return p.ifStatement()
	}
	if p.match(token.While) {
		return p.whileStatement()
	}
	if p.match(token.Print) {
		return p.printStatement()
	}
	if p.match(token.LeftBrace) {
		var statements []ast.Stmt
		statements, err = p.block()
		if err != nil {
			return
		}
		return &ast.BlockStmt{
			Statements: statements,
		}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (stmt ast.Stmt, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after if statement: %w", err)
		return
	}

	var initializer ast.Stmt
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return
		}
	}

	var condition ast.Expr
	if !p.check(token.Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return
		}
	}
	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after condition: %w", err)
		return
	}

	var increment ast.Expr
	if !p.check(token.RightParen) {
		increment, err = p.expression()
		if err != nil {
			return
		}
	}
	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after for loop conditions: %w", err)
		return
	}

	body, err := p.statement()
	if err != nil {
		return
	}
	// We begin the de-sugaring:
	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				body,
				&ast.ExpressionStmt{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &ast.LiteralExpr{Value: true}
	}
	body = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (stmt ast.Stmt, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after if statement: %w", err)
		return
	}
	expr, err := p.expression()
	if err != nil {
		return
	}

	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after if condition: %w", err)
		return
	}

	thenBranch, err := p.statement()
	if err != nil {
		return
	}
	var elseBranch ast.Stmt
	if p.match(token.Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return
		}
	}

	return &ast.IfStmt{
		Condition:  expr,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) whileStatement() (stmt ast.Stmt, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after if statement: %w", err)
		return
	}

	expr, err := p.expression()
	if err != nil {
		return
	}
	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after if condition: %w", err)
		return
	}

	body, err := p.statement()
	if err != nil {
		return
	}

	return &ast.WhileStmt{
		Condition: expr,
		Body:      body,
	}, nil
}

func (p *Parser) block() (stmts []ast.Stmt, err error) {
	for !p.check(token.RightBrace) && !p.isAtEnd() {
		var stmt ast.Stmt
		stmt, err = p.declaration()
		if err != nil {
			return
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.consume(token.RightBrace)
	if err != nil {
		err = fmt.Errorf("expected '}' after block: %w", err)
		return
	}
	return
}

func (p *Parser) printStatement() (stmt ast.Stmt, err error) {
	value, err := p.expression()
	if err != nil {
		return
	}
	_, err = p.consume(token.Semicolon)
	if err != nil {
		return
	}
	return &ast.PrintStmt{
		Expression: value,
	}, nil
}

func (p *Parser) expressionStatement() (stmt ast.Stmt, err error) {
	expr, err := p.expression()
	if err != nil {
		return
	}
	_, err = p.consume(token.Semicolon)
	if err != nil {
		return
	}
	return &ast.ExpressionStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) expression() (expr ast.Expr, err error) {
	return p.assignment()
}
func (p *Parser) assignment() (expr ast.Expr, err error) {
	expr, err = p.or()
	if err != nil {
		return
	}
	if !p.match(token.Equal) {
		return
	}

	equals := p.previous()
	var rvalue ast.Expr
	rvalue, err = p.assignment()
	if err != nil {
		err = fmt.Errorf("invalid rvalue %w", err)
		return
	}

	if varExpr, ok := expr.(*ast.VariableExpr); ok {
		expr = &ast.AssignExpr{
			Name:  varExpr.Name,
			Value: rvalue,
		}
		return
	}

	err = fmt.Errorf("Invalid assignment target for equals token, at line: %d", equals.Line)
	return

}

func (p *Parser) or() (expr ast.Expr, err error) {
	expr, err = p.and()
	if err != nil {
		return
	}
	for p.match(token.Or) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.and()
		if err != nil {
			return
		}
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return
}

func (p *Parser) and() (expr ast.Expr, err error) {
	expr, err = p.equality()
	if err != nil {
		return
	}
	for p.match(token.And) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.equality()
		if err != nil {
			return
		}
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return
}

func (p *Parser) equality() (expr ast.Expr, err error) {
	expr, err = p.comparison()
	if err != nil {
		return
	}
	for p.match(token.EqualEqual, token.BangEqual) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.comparison()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return
}

func (p *Parser) comparison() (expr ast.Expr, err error) {
	expr, err = p.term()
	if err != nil {
		return
	}
	for p.match(token.Less, token.LessEqual, token.Greater, token.GreaterEqual) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.term()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return
}

func (p *Parser) term() (expr ast.Expr, err error) {
	expr, err = p.factor()
	if err != nil {
		return
	}
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.factor()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return
}

func (p *Parser) factor() (expr ast.Expr, err error) {
	expr, err = p.unary()
	if err != nil {
		return
	}

	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.unary()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return
}

// unary → ( "!" | "-" ) unary | call ;
func (p *Parser) unary() (expr ast.Expr, err error) {
	if p.match(token.Slash, token.Star) {
		operator := p.previous()
		var right ast.Expr
		right, err = p.unary()
		if err != nil {
			return
		}
		expr = &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}
		return
	}
	return p.call()
}

// call      → primary ( "(" arguments? ")" )* ;
// arguments → expression ( "," expression )* ;
func (p *Parser) call() (expr ast.Expr, err error) {
	expr, err = p.primary()
	if err != nil {
		return
	}
	for {
		if p.match(token.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return
			}
		} else {
			break
		}
	}
	return
}

func (p *Parser) finishCall(incoming ast.Expr) (expr ast.Expr, err error) {
	var args []ast.Expr
	if !p.check(token.RightParen) {
		for {
			if len(args) >= maxFuncArgCounts {
				// Having a maximum number of arguments will simplify our bytecode interpreter
				// in Part III. We want our two interpreters to be compatible with each other,
				// even in weird corner cases like this, so we’ll add the same limit
				err = fmt.Errorf("Can't have more than 255 arguments.")
				return
			}
			var arg ast.Expr
			arg, err = p.expression()
			if err != nil {
				return
			}
			args = append(args, arg)
			if !p.match(token.Comma) {
				break
			}
		}
	}
	paren, err := p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after arguments: %w", err)
		return
	}
	return &ast.CallExpr{
		Paren:  paren,
		Callee: incoming,
		Args:   args,
	}, nil
}

func (p *Parser) primary() (expr ast.Expr, err error) {
	if p.match(token.False) {
		expr = &ast.LiteralExpr{
			Value: false,
		}
		return
	}
	if p.match(token.True) {
		expr = &ast.LiteralExpr{
			Value: true,
		}
		return
	}
	if p.match(token.Nil) {
		expr = &ast.LiteralExpr{
			Value: nil,
		}
		return
	}
	if p.match(token.Number, token.String) {
		expr = &ast.LiteralExpr{
			Value: p.previous().Literal,
		}
		return
	}
	if p.match(token.Identifier) {
		expr = &ast.VariableExpr{
			Name: p.previous(),
		}
		return
	}

	if p.match(token.LeftParen) {
		var groupingExpr ast.Expr
		groupingExpr, err = p.expression()
		if err != nil {
			return
		}
		_, err = p.consume(token.RightParen)
		if err != nil {
			return
		}
		expr = &ast.GroupingExpr{
			Expression: groupingExpr,
		}
		return
	}
	actual := p.peek()
	err = fmt.Errorf("unexpected token type: %s, at line: %d", actual.TokenType, actual.Line)
	return
}

// ----------------------------------------------------------------------------
// Public API

func (p *Parser) Parse() (statements []ast.Stmt, err error) {
	statements = make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		statement, statementErr := p.declaration()
		if statementErr != nil {
			err = ParserError{
				err: statementErr,
			}
			return
		}
		statements = append(statements, statement)
	}
	return
}
