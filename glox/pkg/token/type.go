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
	// And
	// Class
	// Else
	// False
	// Fun
	// For
	// If
	// Nil
	// Or
	// Print
	// Return
	// Super
	// This
	// True
	// Var
	// While
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

	// case Number:
	// 	return ""
	// case And:
	// 	return ""
	// case Class:
	// 	return ""
	// case Else:
	// 	return ""
	// case False:
	// 	return ""
	// case Fun:
	// 	return ""
	// case For:
	// 	return ""
	// case If:
	// 	return ""
	// case Nil:
	// 	return ""
	// case Or:
	// 	return ""
	// case Print:
	// 	return ""
	// case Return:
	// 	return ""
	// case Super:
	// 	return ""
	// case This:
	// 	return ""
	// case True:
	// 	return ""
	// case Var:
	// 	return ""
	// case While:
	// 	return ""
	case Eof:
		return "EOF"
	default:
		// Is there a better pattern here?
		return "invalid type!"
	}
}
