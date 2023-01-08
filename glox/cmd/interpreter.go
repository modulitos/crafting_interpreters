package main

import (
	"fmt"
	"os"

	"github.com/modulitos/glox/pkg/lox"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := lox.RunFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
	} else {
		err := lox.RunPrompt()
		if err != nil {
			err = fmt.Errorf("exiting due to error: %w", err)
		}
		os.Exit(65)
	}
}
