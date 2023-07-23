package main

import (
	"fmt"
	"os"

	"github.com/Jamlee977/CustomLanguage/parser"
	"github.com/Jamlee977/CustomLanguage/runtimelang"
)

func repl() {
    parser := parser.NewParser()
    fmt.Println("REPL v0.2")
    for {
        fmt.Print("> ")
        var input string
        fmt.Scanln(&input)
        if input == "exit" {
            os.Exit(0)
        }

        program := parser.ProduceAST(input)

        runtimeValue, err := runtimelang.Evaluate(&program)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        fmt.Println(runtimeValue.Get(), runtimeValue.Type())
    }
}

func main() {
    repl()
}
