// Golang's parser:
// https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/parser/parser.go#524
package parser

import (
	"fmt"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
)

type parser struct {
	tokens  []*token.Token
	current int
}

// ----------------------------------------------------------------------------
// Parsing support

func (p *parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == token.Eof
}

func (p *parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	} else {
		return tokenType == p.tokens[p.current].TokenType
	}
}

func (p *parser) match(types ...token.Type) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) consume(tokenType token.Type) (*token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	} else {
		actual := p.peek()
		return nil, fmt.Errorf("wanted token type: %s, got: %s, at line: %d", tokenType, actual, actual.Line)
	}
}

// ----------------------------------------------------------------------------
// Types

func (p *parser) expression() (expr ast.Expr, err error) {
	return p.equality()
}

func (p *parser) equality() (expr ast.Expr, err error) {
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

func (p *parser) comparison() (expr ast.Expr, err error) {
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

func (p *parser) term() (expr ast.Expr, err error) {
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

func (p *parser) factor() (expr ast.Expr, err error) {
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

func (p *parser) unary() (expr ast.Expr, err error) {
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

func (p *parser) primary() (expr ast.Expr, err error) {
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

func (p *parser) parse() ast.Expr {
	expr, err := p.expression()
	if err != nil {
		fmt.Println("err: %w", err)
		return nil
	}
	return expr
}
