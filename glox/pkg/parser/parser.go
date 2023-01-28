// Golang's parser:
// https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/parser/parser.go#524
package parser

import (
	"fmt"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

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

func (p *Parser) declaration() (stmt ast.Stmt, err error) {
	defer func() {
		if err != nil {
			p.synchronize()
			return
		}
	}()

	if p.match(token.Var) {
		return p.varDeclaration()
	}
	return p.statement()
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
	if p.match(token.If) {
		return p.ifStatement()
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
	expr, err = p.equality()
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
	return p.primary()
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
