package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/modulitos/glox/pkg/interpreter"
	"github.com/modulitos/glox/pkg/parser"
	"github.com/modulitos/glox/pkg/scanner"
)

func run(source []byte, interpreterInstance *interpreter.Interpreter) error {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		err = fmt.Errorf("Scanning tokens: %w", err)
		return err
	}

	parser := parser.Parser{Tokens: tokens}
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	resolver := interpreter.NewResolver(interpreterInstance)
	err = resolver.ResolveStmts(statements)
	if err != nil {
		return err
	}
	return interpreterInstance.Interpret(statements)
}

func RunFile(file string) error {
	fmt.Printf("running file: %s\n", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Reading script file: %w", err)
	}
	interpreter := interpreter.NewInterpreter(os.Stdout)
	return run(bytes, interpreter)
}

func RunPrompt() (err error) {
	fmt.Println("starting up lox version 0.0.0")
	scanner := bufio.NewScanner(os.Stdin)
	interpreter := interpreter.NewInterpreter(os.Stdout)
	fmt.Print("> ")
	for scanner.Scan() {
		promptErr := run(scanner.Bytes(), interpreter)
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
