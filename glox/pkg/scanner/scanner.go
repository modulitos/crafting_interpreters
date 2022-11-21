package scanner

import (
	"github.com/modulitos/glox/pkg/token"
)

type scanner struct {
	source []byte
	tokens []*token.Token

	// Points to the first character in the lexeme being scanned.
	start int

	// Points at the character currently being considered.
	current int

	// Tracks what source line current is on so we can produce tokens that know their location.
	line int
}

func NewScanner(source []byte) scanner {
	return scanner{
		source:  source,
		tokens:  []*token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) scanToken() {

}

func (s *scanner) ScanTokens() ([]*token.Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewEofToken(s.line))
	return s.tokens, nil
}
