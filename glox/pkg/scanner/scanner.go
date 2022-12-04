package scanner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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
			if s.isDigit(c) {
				s.number()
			} else if s.isAlpha(c) {
				s.identifier()
			} else {
				s.errors = append(s.errors, fmt.Errorf("Unexpected character: %c on line: %d", c, s.line))
			}
			return
		}
	}
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() rune {
	c, size := utf8.DecodeRune(s.source[s.current:])
	s.current += size
	return rune(c)
}

func (s *scanner) isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func (s *scanner) isAlpha(c rune) bool {
	return 'a' <= c && c <= 'z' ||
		'A' <= c && c <= 'Z' ||
		c == '_'
}

func (s *scanner) isAlphaNumeric(c rune) bool {
	return s.isDigit(c) || s.isAlpha(c)
}

func (s *scanner) addToken(t token.Type, literal interface{}) {
	s.tokens = append(s.tokens,
		&token.Token{
			TokenType: t,
			Lexeme:    string(s.source[s.start:s.current]),
			Line:      s.line,
			Literal:   literal,
		})
}

func (s *scanner) addSimpleToken(t token.Type) {
	s.addToken(t, nil)
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	} else if actual, _ := utf8.DecodeRune(s.source[s.current:]); actual != expected {
		return false
	} else {
		s.current++
		return true
	}
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return utf8.RuneError
	} else {
		r, _ := utf8.DecodeRune(s.source[s.current:])
		return r
	}
}

func (s *scanner) peekNext() rune {
	if s.isAtEnd() {
		return utf8.RuneError
	} else {
		_, size := utf8.DecodeRune(s.source[s.current:])
		if s.current+size >= len(s.source) {
			return utf8.RuneError
		}
		r, _ := utf8.DecodeRune(s.source[s.current+size:])
		return r
	}
}

func (s *scanner) string() (err error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		} else {
			s.advance()
		}
	}

	if s.isAtEnd() {
		err = fmt.Errorf("Unterminated string at line: %d", s.line)
		return
	}

	s.advance()
	s.addToken(token.String, string(s.source[s.start+1:s.current-1]))
	return
}

func (s *scanner) number() (err error) {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	str := string(s.source[s.start:s.current])
	literal, err := strconv.ParseFloat(str, 64)

	if err != nil {
		err = fmt.Errorf("Failed to parse number into string: %s, at line: %d", str, s.line)
		return
	}

	s.addToken(token.Number, literal)
	return
}

func (s *scanner) identifier() (err error) {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	s.addToken(token.Identifier, nil)
	return
}
