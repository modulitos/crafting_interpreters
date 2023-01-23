package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/modulitos/glox/pkg/interpreter"
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

	parser := parser.Parser{Tokens: tokens}
	statements, err := parser.Parse()
	if err != nil {
		return
	}

	interpreter := interpreter.NewInterpreter(os.Stdout)
	err = interpreter.Interpret(statements)

	return
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
