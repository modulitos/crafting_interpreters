package interpreter

import (
	"fmt"
	"time"

	"github.com/modulitos/glox/pkg/ast"
)

// ////////////////////////////////////////////////////////////////////////////
// Native Functions
// ////////////////////////////////////////////////////////////////////////////
type nativeFuncClock struct{}

func (f *nativeFuncClock) String() string {
	return "<native fn>"
}

func (f *nativeFuncClock) arity() int {
	return 0
}
func (f *nativeFuncClock) call(interpreter *Interpreter, args []interface{}) (result interface{}, err error) {
	result = float64(time.Now().UnixMilli()) / 1000.0
	return
}

//////////////////////////////////////////////////////////////////////////////
// Lox Callable Function
//////////////////////////////////////////////////////////////////////////////

// We don’t want the runtime phase of the interpreter to bleed into the front
// end’s syntax classes so we don’t want ast.FunctionStmt itself to implement that.
// Instead, we wrap it in a new class.
type loxFunction struct {
	declaration *ast.FunctionStmt
}

func (f *loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *loxFunction) arity() int {
	return len(f.declaration.Params)
}

func (f *loxFunction) call(interpreter *Interpreter, args []interface{}) (result interface{}, err error) {
	// To support recursion, we create a new environment at each _call_, not at
	// the function declaration.
	environment := newEnvironment(interpreter.globals)
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.define(f.declaration.Params[i].Lexeme, args[i])
	}
	err = interpreter.executeBlock(f.declaration.Body, environment)
	return
}
