package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/modulitos/glox/pkg/scanner"
)

func run(source []byte) error {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		return fmt.Errorf("Scanning tokens: %w", err)
	}
	fmt.Printf("tokens:\n")
	for _, token := range tokens {
		fmt.Println(token)
	}
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

func RunPrompt() error {
	// fmt.Println("running prompt!")
	// os.ReadLin
	// return fmt.Errorf("Not implemented!")
	fmt.Println("starting up lox version 0.0.0")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		run(scanner.Bytes())
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}
	return nil
}
