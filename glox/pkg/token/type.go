package token

type tokenType int

const (
	// Single-character tokens.
	LeftParen = tokenType(iota)
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
	Bang_equal
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// // Literals.
	// Identifier
	// String
	// Number

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

func (t tokenType) toString() string {
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
	case Bang_equal:
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
	// case Identifier:
	// return "="
	// case String:
	// 	return ""
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
