package ast

import (
	"testing"

	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestNaiveExprPrinter(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name:     "literal",
			expr:     &LiteralExpr{Value: 123.489},
			expected: "123.489",
		},
		{
			name:     "nil",
			expr:     &LiteralExpr{Value: nil},
			expected: "nil",
		},
		{
			name: "unary",
			expr: &UnaryExpr{
				Operator: &token.Token{
					TokenType: token.Minus,
					Lexeme:    "-",
				},
				Right: &LiteralExpr{Value: 456},
			},
			expected: "(- 456)",
		},
		{
			name: "binary",
			expr: &BinaryExpr{
				Left: &LiteralExpr{Value: 123},
				Operator: &token.Token{
					TokenType: token.Plus,
					Lexeme:    "+",
				},
				Right: &LiteralExpr{Value: 456},
			},
			expected: "(+ 123 456)",
		},
		{
			name: "grouping",
			expr: &GroupingExpr{
				Expression: &LiteralExpr{Value: 123},
			},
			expected: "(group 123)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			p := AstPrint{}

			// When:
			actual, err := p.Print(&tc.expr)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}

			// Then:
			assert.Equal(t, tc.expected, actual)
		})
	}
}
