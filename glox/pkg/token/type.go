package token

type Type int

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

func (t Type) String() string {
	switch t {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case Comma:
		return ","
	case Dot:
		return "."
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Semicolon:
		return ";"
	case Slash:
		return "/"
	case Star:
		return "*"
	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case EqualEqual:
		return "=="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Less:
		return "<"
	case LessEqual:
		return "<="
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"
	case Eof:
		return "EOF"
	default:
		if str, exists := reverseKeywords[t]; exists == true {
			return str
		} else {
			// Is there a better pattern here?
			return "invalid type!"
		}
	}
}

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

var reverseKeywords = map[Type]string{
	And:    "and",
	Class:  "class",
	Else:   "else",
	False:  "false",
	Fun:    "fun",
	For:    "for",
	If:     "if",
	Nil:    "nil",
	Or:     "or",
	Print:  "print",
	Return: "return",
	Super:  "super",
	This:   "this",
	True:   "true",
	Var:    "var",
	While:  "while",
}
