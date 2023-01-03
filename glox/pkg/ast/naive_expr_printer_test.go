package ast

import (
	"testing"

	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestNaiveExprPrinter(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name:     "literal",
			expr:     &LiteralExpr{value: 123.489},
			expected: "123.489",
		},
		{
			name:     "nil",
			expr:     &LiteralExpr{value: nil},
			expected: "nil",
		},
		{
			name: "unary",
			expr: &UnaryExpr{
				operator: &token.Token{
					TokenType: token.Minus,
					Lexeme:    "-",
				},
				right: &LiteralExpr{value: 456},
			},
			expected: "(- 456)",
		},
		{
			name: "binary",
			expr: &BinaryExpr{
				left: &LiteralExpr{value: 123},
				operator: &token.Token{
					TokenType: token.Plus,
					Lexeme:    "+",
				},
				right: &LiteralExpr{value: 456},
			},
			expected: "(+ 123 456)",
		},
		{
			name: "grouping",
			expr: &GroupingExpr{
				expression: &LiteralExpr{value: 123},
			},
			expected: "(group 123)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			p := AstPrint{}

			// When:
			actual, err := p.print(&tc.expr)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}

			// Then:
			require.Equal(t, tc.expected, actual)
		})
	}
}
