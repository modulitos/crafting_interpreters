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
			name:       "unexpected characters",
			source:     "(+)^ {.}^",
			wantErr:    fmt.Errorf("Unexpected character: ^ on line: 1\nUnexpected character: ^ on line: 1\n"),
			wantTokens: nil,
		},
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
