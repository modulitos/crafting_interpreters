package interpreter

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestInterpreterEvaluate(t *testing.T) {
	tests := []struct {
		name        string
		stmts       []ast.Stmt
		expected    string
		expectedErr bool
	}{
		{
			name: "float",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.LiteralExpr{Value: 123.489},
				},
			},
			expected: "123.489\n",
		},
		{
			name: "float without decimals",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.LiteralExpr{Value: 123.0},
				},
			},
			expected: "123\n",
		},
		{
			name: "addition",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 1.0},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 2.0},
					},
				},
			},
			expected: "3\n",
		},
		{
			name: "addition with grouping multiplication",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
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
				},
			},
			expected: "21.1\n",
		},
		{
			name: "addition of string and number",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: "asdf"},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 1.0},
					},
				},
			},
			expected: "asdf1\n",
		},
		{
			name: "addition of number and string",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 1.0},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: "asdf"},
					},
				},
			},
			expected: "1asdf\n",
		},
		{
			name: "addition of number and string",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 1.0},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: true},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "divide by 0",
			stmts: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 1.0},
						Operator: &token.Token{
							TokenType: token.Slash,
							Lexeme:    "/",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 0.0},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "declare a variable, then use it",
			stmts: []ast.Stmt{
				&ast.VarStmt{
					Name: &token.Token{
						TokenType: token.Identifier,
						Lexeme:    "foo",
						Line:      1,
					},
					Initializer: &ast.LiteralExpr{
						Value: 2.0,
					},
				},
				&ast.PrintStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.VariableExpr{
							Name: &token.Token{
								TokenType: token.Identifier,
								Lexeme:    "foo",
								Line:      1,
							},
						},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 2.0},
					},
				},
			},
			expected: "4\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			buf := new(bytes.Buffer)
			interpreter := NewInterpreter(buf)

			// When:
			err := interpreter.Interpret(tc.stmts)
			if (err != nil) != tc.expectedErr {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\nexpectedErr:\n%v\n", tc.name, err, tc.expectedErr)

				return
			}

			b, err := ioutil.ReadAll(buf)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}
			actual := string(b)

			// Then:
			assert.Equal(t, tc.expected, actual)
		})
	}
}
