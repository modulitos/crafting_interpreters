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
`))
}

func (g *generator) linebreak() {
	g.buf.WriteByte('\n')
}

func (g *generator) writeTypes(types []string) {
	// TODO: parse types in a GrammarType struct so it's easier

	// write visitor
	g.linebreak()
	fmt.Fprintf(&g.buf, "type ExprVisitor interface {")
	g.linebreak()
	for _, typestr := range types {
		name := strings.TrimSpace(strings.Split(typestr, ":")[0])
		fmt.Fprintf(&g.buf, "Visit%s(e *%sExpr) (result interface{}, err error)", name, name)
		g.linebreak()
	}
	g.buf.Write([]byte("}\n"))

	for _, typestr := range types {
		g.linebreak()
		name := strings.TrimSpace(strings.Split(typestr, ":")[0])
		fields := strings.Split(strings.Split(typestr, ":")[1], ",")

		// Define the Expr struct:
		fmt.Fprintf(&g.buf, "type %sExpr struct {", name)
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
		fmt.Fprintf(&g.buf, "func (e *%sExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {", name)
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
		"Binary : left Expr, operator *token.Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value interface{}",
		"Unary : operator *token.Token, right Expr",
	})

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
