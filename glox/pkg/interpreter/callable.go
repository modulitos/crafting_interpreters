package interpreter

type Callable interface {
	call(interpreter *Interpreter, args []interface{}) (result interface{}, err error)
	arity() int
	String() string
}
