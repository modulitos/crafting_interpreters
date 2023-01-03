package token

// This is how tokens are implemented in golang:
// https://cs.opensource.google/go/go/+/master:src/go/token/token.go
type Type int

// https://cs.opensource.google/go/go/+/master:src/go/ast/ast.go

const (
	// Single-character tokens.
	LeftParen = Type(iota)
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// // Literals.
	Identifier
	String
	Number

	// // Keywords.
	Eof
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
)

var Keywords = map[string]Type{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}
