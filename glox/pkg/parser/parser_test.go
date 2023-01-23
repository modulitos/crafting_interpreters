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
