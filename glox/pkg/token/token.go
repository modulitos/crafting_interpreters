package token

import "fmt"

type Token struct {
	TokenType Type
	// The lexemes are only the raw substrings of the source code.
	Lexeme string
	Line   int

	// Literals are numbers and strings and the like. Since the scanner has to
	// walk each character in the Literal to correctly identify it, it can also
	// convert that textual representation of a value to the living runtime
	// object that will be used by the interpreter later.
	Literal interface{}
}

func NewEofToken(line int) *Token {
	return &Token{
		TokenType: Eof,
		Line:      line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("type: %s with lexeme: %q with literal: %s", t.TokenType.String(), t.Lexeme, t.Literal)
}
