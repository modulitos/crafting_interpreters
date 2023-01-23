package interpreter

import (
	"fmt"

	"github.com/modulitos/glox/pkg/token"
)

type environment struct {
	values map[string]interface{}
	parent *environment
}

func newEnvironment(parent *environment) *environment {
	return &environment{
		values: make(map[string]interface{}),
		parent: parent,
	}
}

func (e *environment) define(name string, value interface{}) {
	// We have made one interesting semantic choice: When we add the key to the
	// map, we don’t check to see if it’s already present.
	//
	// A variable statement doesn’t just define a new variable, it can also be
	// used to redefine an existing variable. We could choose to make this an
	// error instead, but that would interact poorly with the REPL.
	e.values[name] = value
}

func (e *environment) get(name *token.Token) (result interface{}, err error) {
	if result, exists := e.values[name.Lexeme]; exists {
		return result, nil
	} else {
		// Since making it a static error makes recursive declarations too
		// difficult, we'll defer the error to runtime. It's OK to refer to a
		// variable before it's defined as long as you don't evaluate the
		// reference.
		return nil, &RuntimeError{
			msg: fmt.Sprintf("Undefined variable: %s.\n", name.Lexeme),
		}
	}
}
