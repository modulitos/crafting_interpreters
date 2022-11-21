package token

import "fmt"

type Token struct {
	tokenType tokenType
	// The lexemes are only the raw substrings of the source code.
	lexeme *string
	line   int

	// Literals are numbers and strings and the like. Since the scanner has to
	// walk each character in the literal to correctly identify it, it can also
	// convert that textual representation of a value to the living runtime
	// object that will be used by the interpreter later.
	literal *string
}

func NewEofToken(line int) *Token {
	return &Token{
		tokenType: Eof,
		lexeme:    nil,
		line:      line,
		literal:   nil,
	}
}

func (t *Token) toString() string {
	return fmt.Sprintf("%s %s %s", t.tokenType.toString(), *t.lexeme, *t.literal)
}
