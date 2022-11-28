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

func (s *scanner) scanToken() {
	for !s.isAtEnd() {
		c := s.advance()
		switch c {
		case ' ', '\r', '\t':
			// do nothing on whitespace chars
			return
		case '\n':
			s.line++
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
			if s.match('/') {
				// Consume all chars until end of line.
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
				// Comments as lexemes are ignored.
			} else {
				s.addSimpleToken(token.Slash)
			}
			return
		case '*':
			s.addSimpleToken(token.Star)
			return
		case '!':
			if s.match('=') {
				s.addSimpleToken(token.BangEqual)
			} else {
				s.addSimpleToken(token.Bang)
			}
			return
		case '>':
			if s.match('=') {
				s.addSimpleToken(token.GreaterEqual)
			} else {
				s.addSimpleToken(token.Greater)
			}
			return
		case '<':
			if s.match('=') {
				s.addSimpleToken(token.LessEqual)
			} else {
				s.addSimpleToken(token.Less)
			}
			return
		case '=':
			if s.match('=') {
				s.addSimpleToken(token.EqualEqual)
			} else {
				s.addSimpleToken(token.Equal)
			}
			return
		case '"':
			s.string()
			return
		default:
			s.errors = append(s.errors, fmt.Errorf("Unexpected character: %c on line: %d", c, s.line))
			return
		}
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

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	} else if s.source[s.current] != expected {
		return false
	} else {
		s.current++
		return true
	}
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		// TODO: Consider using a rune here instead?
		//
		// byte is basically a u8. So we return 0, which is the default byte:
		return 0
	} else {
		return s.source[s.current]
	}
}

func (s *scanner) string() (err error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.advance()
		}
	}

	if s.isAtEnd() {
		err = fmt.Errorf("Unterminated string at line: %d", s.line)
		return
	}

	s.advance()
	substring := string(s.source[s.start+1 : s.current-1])
	tokenType := token.String
	s.addToken(&tokenType, &substring)
	return
}
