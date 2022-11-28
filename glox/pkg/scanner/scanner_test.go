package scanner

import (
	"fmt"
	"testing"

	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/require"
)

func simpleToken(tokenType token.Type, line int) *token.Token {
	// This is a simplification! We may need to customize the lexeme.
	tokenString := tokenType.String()
	return &token.Token{
		TokenType: tokenType,
		Lexeme:    &tokenString,
		Literal:   nil,
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
				simpleToken(token.LeftParen, 1),
				simpleToken(token.Plus, 1),
				simpleToken(token.RightParen, 1),
				simpleToken(token.LeftBrace, 1),
				simpleToken(token.Dot, 1),
				simpleToken(token.RightBrace, 1),
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
				simpleToken(token.LessEqual, 1),
				simpleToken(token.EqualEqual, 1),
				simpleToken(token.BangEqual, 1),
				simpleToken(token.Bang, 1),
				token.NewEofToken(1),
			},
		},
		{
			name:    "comment",
			source:  "!\n!!// this is a comment \n() // some other comment",
			wantErr: nil,
			wantTokens: []*token.Token{
				simpleToken(token.Bang, 1),
				simpleToken(token.Bang, 2),
				simpleToken(token.Bang, 2),
				simpleToken(token.LeftParen, 3),
				simpleToken(token.RightParen, 3),
				token.NewEofToken(3),
			},
		},
		// {
		// 	name:    "string",
		// 	source:  "\"hello world\"",
		// 	wantErr: nil,
		// 	wantTokens: []*token.Token{
		// 		&token.Token{
		// 			TokenType: token.String,
		// 			Lexeme:    `""`,
		// 			Literal:   &"",
		// 			Line:      1,
		// 		},
		// 		token.NewEofToken(1),
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScanner([]byte(tc.source))
			gotTokens, err := s.ScanTokens()

			// Assert that the errors are either both non-nil, or have the same error message.
			if err != tc.wantErr && err.Error() != tc.wantErr.Error() {
				t.Errorf("ScanTokens() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			require.Equal(t, tc.wantTokens, gotTokens)
		})
	}
}
