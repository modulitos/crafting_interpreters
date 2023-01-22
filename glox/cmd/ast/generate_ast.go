package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
)

type generator struct {
	buf bytes.Buffer
}

func (g *generator) writeHeader() {
	g.buf.Write([]byte(
		`// Code generated by generate_ast. DO NOT EDIT.
// Eg of Go's AST: https://go.googlesource.com/go/+/38cfb3be9d486833456276777155980d1ec0823e/src/go/ast/ast.go#1

package ast

import (
	"github.com/modulitos/glox/pkg/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) (result interface{}, err error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}
`))
}

func (g *generator) linebreak() {
	g.buf.WriteByte('\n')
}

type exprType int

const (
	statement = exprType(iota)
	expression
)

func (g *generator) writeTypes(types []string, kind exprType) {
	var exprRepr string
	var return_type string
	if kind == statement {
		exprRepr = "Stmt"
		return_type = fmt.Sprintf("error")
	} else {
		exprRepr = "Expr"
		return_type = fmt.Sprintf("(result interface{}, err error)")
	}

	// if kind == statement {
	// } else {
	// }

	// write visitor
	g.linebreak()
	fmt.Fprintf(&g.buf, "type %sVisitor interface {", exprRepr)
	g.linebreak()
	for _, typestr := range types {
		name := strings.TrimSpace(strings.Split(typestr, ":")[0])
		fmt.Fprintf(&g.buf, "Visit%s(e *%s%s) %s", name, name, exprRepr, return_type)
		g.linebreak()
	}
	g.buf.Write([]byte("}\n"))

	for _, typestr := range types {
		g.linebreak()
		name := strings.TrimSpace(strings.Split(typestr, ":")[0])
		fields := strings.Split(strings.Split(typestr, ":")[1], ",")

		// Define the Expr struct:
		fmt.Fprintf(&g.buf, "type %s%s struct {", name, exprRepr)
		g.linebreak()
		for _, field := range fields {
			field := strings.TrimSpace(field)
			name := strings.Split(field, " ")[0]
			fieldType := strings.Split(field, " ")[1]
			fmt.Fprintf(&g.buf, `%s %s`, name, fieldType)
			g.linebreak()
		}
		g.buf.Write([]byte("}"))
		g.linebreak()

		// implement the Accept method:
		fmt.Fprintf(&g.buf, "func (e *%s%s) Accept(visitor %sVisitor) %s {", name, exprRepr, exprRepr, return_type)
		g.linebreak()
		fmt.Fprintf(&g.buf, "return visitor.Visit%s(e)", name)
		g.linebreak()
		g.buf.Write([]byte("}"))
		g.linebreak()
	}
}

func (g *generator) format() (err error) {
	formatted, err := format.Source(g.buf.Bytes())
	if err != nil {
		err = fmt.Errorf("Formatting code with gofmt: %w\n\ncode: %s", err, g.buf.String())
		return
	}
	g.buf.Reset()
	g.buf.Write(formatted)
	return
}

func (g *generator) writeTo(writer io.Writer) {
	io.Copy(writer, &g.buf)
}

func main() {
	var output = flag.String("o", "pkg/ast/ast_generated.go", "Usage: go run generate_ast.go -o <output_file>")
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			fmt.Println("fatal: " + err.Error())
			os.Exit(65)
		}
	}()

	fmt.Println("starting!")

	generator := generator{}
	generator.writeHeader()

	generator.writeTypes([]string{
		// Sometimes assign is a statement, like in Go:
		// "Assign : Token name, Expr value",
		"Binary : Left Expr, Operator *token.Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal : Value interface{}",
		"Unary : Operator *token.Token, Right Expr",
	}, expression)

	generator.writeTypes([]string{
		"Expression : Expression Expr",
		"Print : Expression Expr",
		// statements:
		// "Block : List<Stmt> statements",
		// todo: declaration statement
	}, statement)

	err = generator.format()
	if err != nil {
		err = fmt.Errorf("Error formatting code: %w", err)
		return
	}

	f, err := os.Create(*output)
	if err != nil {
		err = fmt.Errorf("creating new file %q: %w", *output, err)
		return
	}
	defer f.Close()
	defer f.Sync() // nolint: errcheck

	generator.writeTo(f)
	fmt.Println("AST written to file:", *output)
}
