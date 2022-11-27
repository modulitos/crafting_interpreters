package scanner

import (
	"errors"
	"fmt"
	"strings"

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

	errors []error
}

func NewScanner(source []byte) scanner {
	return scanner{
		source:  source,
		tokens:  []*token.Token{},
		start:   0,
		current: 0,
		line:    1,
		errors:  []error{},
	}
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *scanner) addToken(t *token.Type, literal *string) {
	// Using an intermediate variable, because otherwise we can't take the
	// address of a string? https://github.com/golang/go/issues/6031
	str := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens,
		&token.Token{
			TokenType: *t,
			Lexeme:    &str,
			Line:      s.line,
			Literal:   literal,
		})
}

func (s *scanner) addSimpleToken(t token.Type) {
	s.addToken(&t, nil)
}

func (s *scanner) scanToken() {
	for !s.isAtEnd() {
		c := s.advance()
		switch c {
		case ' ', '\r', '\t':
			// do nothing on whitespace chars
			return
		case '(':
			s.addSimpleToken(token.LeftParen)
			return
		case ')':
			s.addSimpleToken(token.RightParen)
			return
		case '{':
			s.addSimpleToken(token.LeftBrace)
			return
		case '}':
			s.addSimpleToken(token.RightBrace)
			return
		case ',':
			s.addSimpleToken(token.Comma)
			return
		case '.':
			s.addSimpleToken(token.Dot)
			return
		case '-':
			s.addSimpleToken(token.Minus)
			return
		case '+':
			s.addSimpleToken(token.Plus)
			return
		case ';':
			s.addSimpleToken(token.Semicolon)
			return
		case '/':
			// TODO: handle comments
			s.addSimpleToken(token.Slash)
			return
		case '*':
			s.addSimpleToken(token.Star)
			return
		case '!':
			s.addSimpleToken(token.Bang)
			return
		case '=':
			s.addSimpleToken(token.Equal)
			return
		default:
			s.errors = append(s.errors, fmt.Errorf("Unexpected character: %c on line: %d", c, s.line))
			return
		}
	}
}

func (s *scanner) ScanTokens() ([]*token.Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewEofToken(s.line))

	if len(s.errors) > 0 {
		builder := strings.Builder{}
		for _, err := range s.errors {
			builder.WriteString(fmt.Sprintf("%v\n", err.Error()))
		}
		return nil, errors.New(builder.String())
	} else {
		return s.tokens, nil
	}
}
