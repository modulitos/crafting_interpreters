package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/modulitos/glox/pkg/ast"
	"github.com/modulitos/glox/pkg/parser"
	"github.com/modulitos/glox/pkg/scanner"
)

func run(source []byte) (err error) {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		err = fmt.Errorf("Scanning tokens: %w", err)
		return
	}
	// fmt.Printf("tokens:\n")

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }
	parser := parser.Parser{Tokens: tokens}
	expression, err := parser.Parse()
	if err != nil {
		// fmt.Printf("parser error! %s\n", err)
		// return fmt.Errorf("Parser error: %w", err)
		return
	}

	astPrinter := ast.AstPrint{}
	fmt.Println(astPrinter.Print(expression))

	return nil
}

func RunFile(file string) error {
	fmt.Printf("running file: %s\n", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Reading script file: %w", err)
	}
	return run(bytes)
}

func RunPrompt() (err error) {
	// fmt.Println("running prompt!")
	// os.ReadLin
	// return fmt.Errorf("Not implemented!")
	fmt.Println("starting up lox version 0.0.0")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		promptErr := run(scanner.Bytes())
		if promptErr != nil {
			fmt.Printf("Error evaluating input: %v\n", promptErr)
		}
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}
	return nil
}
