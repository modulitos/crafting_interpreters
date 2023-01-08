package scanner

import (
	"fmt"
	"testing"

	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/assert"
)

func simpleToken(tokenType token.Type, line int, lexeme string) *token.Token {
	return &token.Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Line:      line,
	}
}

func TestScanner_ScanTokens(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantErr    error
		wantTokens []*token.Token
	}{
		{
			name:    "empty",
			source:  "",
			wantErr: nil,
			wantTokens: []*token.Token{
				token.NewEofToken(1),
			},
		},
		{
			name:    "basic single char",
			source:  "(+) {.}",
			wantErr: nil,
			wantTokens: []*token.Token{
				simpleToken(token.LeftParen, 1, "("),
				simpleToken(token.Plus, 1, "+"),
				simpleToken(token.RightParen, 1, ")"),
				simpleToken(token.LeftBrace, 1, "{"),
				simpleToken(token.Dot, 1, "."),
				simpleToken(token.RightBrace, 1, "}"),
				token.NewEofToken(1),
			},
		},
		{
			name:       "multiple unexpected characters",
			source:     "(+)^ \n {.}^",
			wantErr:    fmt.Errorf("Unexpected character: ^ on line: 1\nUnexpected character: ^ on line: 2\n"),
			wantTokens: nil,
		},
		{
			name:    "multi-byte operators",
			source:  "<= == != !",
			wantErr: nil,
			wantTokens: []*token.Token{
				simpleToken(token.LessEqual, 1, "<="),
				simpleToken(token.EqualEqual, 1, "=="),
				simpleToken(token.BangEqual, 1, "!="),
				simpleToken(token.Bang, 1, "!"),
				token.NewEofToken(1),
			},
		},
		{
			name:    "comment",
			source:  "!\n!!// this is a comment \n() // some other comment",
			wantErr: nil,
			wantTokens: []*token.Token{
				simpleToken(token.Bang, 1, "!"),
				simpleToken(token.Bang, 2, "!"),
				simpleToken(token.Bang, 2, "!"),
				simpleToken(token.LeftParen, 3, "("),
				simpleToken(token.RightParen, 3, ")"),
				token.NewEofToken(3),
			},
		},
		{
			name:    "string",
			source:  "\"hello world\"",
			wantErr: nil,
			wantTokens: []*token.Token{
				{
					TokenType: token.String,
					Lexeme:    "\"hello world\"",
					Literal:   "hello world",
					Line:      1,
				},
				token.NewEofToken(1),
			},
		},
		{
			name:    "number",
			source:  " 123",
			wantErr: nil,
			wantTokens: []*token.Token{
				{
					TokenType: token.Number,
					Lexeme:    "123",
					Literal:   123.0,
					Line:      1,
				},
				token.NewEofToken(1),
			},
		},
		{
			name:    "number with decimal",
			source:  "2345.2342 ",
			wantErr: nil,
			wantTokens: []*token.Token{
				{
					TokenType: token.Number,
					Lexeme:    "2345.2342",
					Literal:   2345.2342,
					Line:      1,
				},
				token.NewEofToken(1),
			},
		},
		{
			name:    "number with method call",
			source:  "2345.foo() ",
			wantErr: nil,
			wantTokens: []*token.Token{
				{
					TokenType: token.Number,
					Lexeme:    "2345",
					Literal:   2345.0,
					Line:      1,
				},
				simpleToken(token.Dot, 1, "."),
				{
					TokenType: token.Identifier,
					Lexeme:    "foo",
					Literal:   nil,
					Line:      1,
				},
				simpleToken(token.LeftParen, 1, "("),
				simpleToken(token.RightParen, 1, ")"),
				token.NewEofToken(1),
			},
		},
		{
			name:    "literal",
			source:  "blah123",
			wantErr: nil,
			wantTokens: []*token.Token{
				{
					TokenType: token.Identifier,
					Lexeme:    "blah123",
					Literal:   nil,
					Line:      1,
				},
				token.NewEofToken(1),
			},
		},
		{
			name:    "key words",
			source:  "var blah123",
			wantErr: nil,
			wantTokens: []*token.Token{
				simpleToken(token.Var, 1, "var"),
				{
					TokenType: token.Identifier,
					Lexeme:    "blah123",
					Literal:   nil,
					Line:      1,
				},
				token.NewEofToken(1),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScanner([]byte(tc.source))
			gotTokens, err := s.ScanTokens()

			// Assert that the errors are either both non-nil, or have the same error message.
			if (err == nil) != (tc.wantErr == nil) {
				t.Errorf("ScanTokens() has an unexpected err:\nerror:\n%v\nwantErr:\n%v\n", err, tc.wantErr)
				return
			}

			if err != tc.wantErr && err.Error() != tc.wantErr.Error() {
				t.Errorf("ScanTokens() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.wantTokens, gotTokens)
		})
	}
}
