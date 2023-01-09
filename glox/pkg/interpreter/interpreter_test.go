package interpreter

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected interface{}
	}{
		{
			name:     "float",
			expr:     &ast.LiteralExpr{Value: 123.489},
			expected: "123.489\n",
		},
		{
			name:     "float without decimals",
			expr:     &ast.LiteralExpr{Value: 123.0},
			expected: "123\n",
		},
		{
			name: "addition",
			expr: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: 1.0},
				Operator: &token.Token{
					TokenType: token.Plus,
					Lexeme:    "+",
					Line:      1,
				},
				Right: &ast.LiteralExpr{Value: 2.0},
			},
			expected: "3\n",
		},
		{
			name: "addition with grouping multiplication",
			expr: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: 1.1},
				Operator: &token.Token{
					TokenType: token.Plus,
					Lexeme:    "+",
					Line:      1,
				},
				Right: &ast.GroupingExpr{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 10.0},
						Operator: &token.Token{
							TokenType: token.Star,
							Lexeme:    "*",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 2.0},
					},
				},
			},
			expected: "21.1\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			buf := new(bytes.Buffer)
			interpreter := NewInterpreter(buf)

			// When:
			err := interpreter.Interpret(tc.expr)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}

			// Then:
			b, err := ioutil.ReadAll(buf)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}
			actual := string(b)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
