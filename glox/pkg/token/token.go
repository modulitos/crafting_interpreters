package token

import "fmt"

type Token struct {
	TokenType Type
	// The lexemes are only the raw substrings of the source code.
	Lexeme *string
	Line   int

	// Literals are numbers and strings and the like. Since the scanner has to
	// walk each character in the Literal to correctly identify it, it can also
	// convert that textual representation of a value to the living runtime
	// object that will be used by the interpreter later.
	Literal *string
}

func NewEofToken(line int) *Token {
	return &Token{
		TokenType: Eof,
		Lexeme:    nil,
		Line:      line,
		Literal:   nil,
	}
}

func (t *Token) String() string {
	lexeme := "nil"
	if t.Lexeme != nil {
		lexeme = *t.Lexeme
	}
	literal := "nil"
	if t.Literal != nil {
		literal = *t.Literal
	}
	return fmt.Sprintf("%s %s %s", t.TokenType.String(), lexeme, literal)
}
