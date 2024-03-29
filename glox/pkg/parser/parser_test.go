package parser

import (
	"testing"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		tokens      []*token.Token
		expected    []ast.Stmt
		expectedErr error
	}{
		{
			name: "simple_binary_expr",
			tokens: []*token.Token{
				{
					TokenType: token.Number,
					Lexeme:    "1",
					Literal:   1,
					Line:      1,
				},
				{
					TokenType: token.Plus,
					Lexeme:    "+",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Number,
					Lexeme:    "2",
					Literal:   2,
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      1,
				},
			},
			expected: []ast.Stmt{
				&ast.ExpressionStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: 1},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: 2},
					},
				},
			},
		},
		{
			name: "binary grouping expr: asdf <= (1+2)",
			tokens: []*token.Token{
				{
					TokenType: token.String,
					Lexeme:    "asdf",
					Literal:   "asdf",
					Line:      1,
				},
				{
					TokenType: token.LessEqual,
					Lexeme:    "<=",
					Line:      1,
				},
				{
					TokenType: token.LeftParen,
					Lexeme:    "(",
					Line:      1,
				},
				{
					TokenType: token.Number,
					Lexeme:    "1",
					Literal:   1,
					Line:      1,
				},
				{
					TokenType: token.Plus,
					Lexeme:    "+",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Number,
					Lexeme:    "2",
					Literal:   2,
					Line:      1,
				},
				{
					TokenType: token.RightParen,
					Lexeme:    ")",
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      1,
				},
			},
			expected: []ast.Stmt{
				&ast.ExpressionStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: "asdf"},
						Operator: &token.Token{
							TokenType: token.LessEqual,
							Lexeme:    "<=",
							Line:      1,
						},
						Right: &ast.GroupingExpr{
							Expression: &ast.BinaryExpr{
								Left: &ast.LiteralExpr{Value: 1},
								Operator: &token.Token{
									TokenType: token.Plus,
									Lexeme:    "+",
									Line:      1,
								},
								Right: &ast.LiteralExpr{Value: 2},
							},
						},
					},
				},
			},
		},
		{
			name: "binary expr: asdf + qwer",
			tokens: []*token.Token{
				{
					TokenType: token.String,
					Lexeme:    "asdf",
					Literal:   "asdf",
					Line:      1,
				},
				{
					TokenType: token.Plus,
					Lexeme:    "+",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.String,
					Lexeme:    "qwer",
					Literal:   "qwer",
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      1,
				},
			},
			expected: []ast.Stmt{
				&ast.ExpressionStmt{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{Value: "asdf"},
						Operator: &token.Token{
							TokenType: token.Plus,
							Lexeme:    "+",
							Line:      1,
						},
						Right: &ast.LiteralExpr{Value: "qwer"},
					},
				},
			},
		},
		{
			name: "print statement",
			tokens: []*token.Token{
				{
					TokenType: token.Print,
					Lexeme:    "print",
					Literal:   "asdf",
					Line:      1,
				},
				{
					TokenType: token.String,
					Lexeme:    "qwer",
					Literal:   "qwer",
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      1,
				},
			},
			expected: []ast.Stmt{
				&ast.PrintStmt{
					Expression: &ast.LiteralExpr{Value: "qwer"},
				},
			},
		},
		{
			name: "variable declaration",
			tokens: []*token.Token{
				{
					TokenType: token.Var,
					Lexeme:    "Var",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Identifier,
					Lexeme:    "qwer",
					Literal:   "qwer",
					Line:      1,
				},
				{
					TokenType: token.Equal,
					Lexeme:    "=",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Number,
					Lexeme:    "42",
					Literal:   42,
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      1,
				},
			},
			expected: []ast.Stmt{
				&ast.VarStmt{
					Name: &token.Token{
						TokenType: token.Identifier,
						Literal:   "qwer",
						Lexeme:    "qwer",
						Line:      1,
					},
					Initializer: &ast.LiteralExpr{Value: 42},
				},
			},
		},
		{
			name: "blocks",
			tokens: []*token.Token{
				{
					TokenType: token.Var,
					Lexeme:    "Var",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Identifier,
					Lexeme:    "foo",
					Literal:   "foo",
					Line:      1,
				},
				{
					TokenType: token.Equal,
					Lexeme:    "=",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.Number,
					Lexeme:    "42",
					Literal:   42,
					Line:      1,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      1,
				},
				{
					TokenType: token.LeftBrace,
					Lexeme:    "{",
					Literal:   nil,
					Line:      2,
				},
				{
					TokenType: token.Var,
					Lexeme:    "Var",
					Literal:   nil,
					Line:      3,
				},
				{
					TokenType: token.Identifier,
					Lexeme:    "foo",
					Literal:   "foo",
					Line:      3,
				},
				{
					TokenType: token.Equal,
					Lexeme:    "=",
					Literal:   nil,
					Line:      3,
				},
				{
					TokenType: token.Number,
					Lexeme:    "42",
					Literal:   42,
					Line:      3,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
					Line:      3,
				},
				{
					TokenType: token.RightBrace,
					Lexeme:    "}",
					Literal:   nil,
					Line:      4,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
					Line:      4,
				},
			},
			expected: []ast.Stmt{
				&ast.VarStmt{
					Name: &token.Token{
						TokenType: token.Identifier,
						Literal:   "foo",
						Lexeme:    "foo",
						Line:      1,
					},
					Initializer: &ast.LiteralExpr{Value: 42},
				},

				&ast.BlockStmt{
					Statements: []ast.Stmt{
						&ast.VarStmt{
							Name: &token.Token{
								TokenType: token.Identifier,
								Literal:   "foo",
								Lexeme:    "foo",
								Line:      3,
							},
							Initializer: &ast.LiteralExpr{Value: 42},
						},
					},
				},
			},
		},
		{

			name: `function: sayHi("Dear", "Reader");`,
			tokens: []*token.Token{
				{
					TokenType: token.Identifier,
					Lexeme:    "sayHi",
					Literal:   "sayHi",
				},
				{
					TokenType: token.LeftParen,
					Lexeme:    "(",
					Literal:   nil,
				},
				{
					TokenType: token.String,
					Lexeme:    "Dear",
					Literal:   "Dear",
				},
				{
					TokenType: token.Comma,
					Lexeme:    ",",
					Literal:   nil,
				},
				{
					TokenType: token.String,
					Lexeme:    "Reader",
					Literal:   "Reader",
				},
				{
					TokenType: token.RightParen,
					Lexeme:    ")",
					Literal:   nil,
				},
				{
					TokenType: token.Semicolon,
					Lexeme:    ";",
					Literal:   nil,
				},
				{
					TokenType: token.Eof,
					Lexeme:    "",
					Literal:   nil,
				},
			},
			expected: []ast.Stmt{
				&ast.ExpressionStmt{
					Expression: &ast.CallExpr{
						Callee: &ast.VariableExpr{
							Name: &token.Token{
								TokenType: token.Identifier,
								Lexeme:    "sayHi",
								Literal:   "sayHi",
							},
						},
						Paren: &token.Token{
							TokenType: token.RightParen,
							Lexeme:    ")",
							Literal:   nil,
						},
						Args: []ast.Expr{
							&ast.LiteralExpr{
								Value: "Dear",
							},
							&ast.LiteralExpr{
								Value: "Reader",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := Parser{
				Tokens: tc.tokens,
			}
			actual, err := parser.Parse()
			if err != nil {
				t.Errorf("Unexpected err:\nerror:\n%v\n", err)
				return
			}
			assert.Equal(t, len(tc.expected), len(actual))
			assert.Equal(t, tc.expected, actual)

		})
	}
}
