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

func newGlobalEnvironment() *environment {
	env := newEnvironment(nil)
	env.define("clock", &nativeFuncClock{})

	return env
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

// The only difference between `assign` and `define` is `assign` isn't allow to
// create a new variable.
func (e *environment) assign(name *token.Token, value interface{}) (err error) {
	if _, exists := e.values[name.Lexeme]; exists {
		e.values[name.Lexeme] = value
	} else {
		if e.parent != nil {
			return e.parent.assign(name, value)
		}

		err = fmt.Errorf("Cannot assign undeclared variable: '%s'.", name.Lexeme)
	}
	return
}

func (e *environment) get(name *token.Token) (result interface{}, err error) {
	if result, exists := e.values[name.Lexeme]; exists {
		return result, nil
	} else {
		if e.parent != nil {
			return e.parent.get(name)
		}

		// Since making it a static error makes recursive declarations too
		// difficult, we'll defer the error to runtime. It's OK to refer to a
		// variable before it's defined as long as you don't evaluate the
		// reference.
		return nil, fmt.Errorf("Undefined variable: %s.\n", name.Lexeme)
	}
}
