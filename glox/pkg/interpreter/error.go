package interpreter

import (
	"fmt"

	"github.com/modulitos/glox/pkg/token"
)

type RuntimeError struct {
	msg   string
	token *token.Token
}

func (e *RuntimeError) Error() string {
	// TODO: don't print the token if it's nil?
	return fmt.Sprintf("Interpreter Runtime Error: %s\nToken: %v\n", e.msg, e.token)
}
