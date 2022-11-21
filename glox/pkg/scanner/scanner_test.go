package scanner

import (
	"testing"

	"github.com/modulitos/glox/pkg/token"
	"github.com/stretchr/testify/require"
)

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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScanner([]byte(tc.source))
			gotTokens, err := s.ScanTokens()
			if err != tc.wantErr {
				t.Errorf("ScanTokens() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			require.Equal(t, tc.wantTokens, gotTokens)
		})
	}
}
