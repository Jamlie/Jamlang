package main

import (
    "fmt"
    "os"

    "github.com/Jamlee977/CustomLanguage/parser"
)

func repl() {
    parser := parser.NewParser()
    fmt.Println("REPL v0.1")
    for {
        fmt.Print("> ")
        var input string
        fmt.Scanln(&input)
        if input == "exit" {
            os.Exit(0)
        }

        program := parser.ProduceAST(input)

        for _, statement := range program.Body {
            fmt.Println(statement.Kind())
            fmt.Println(statement.ToString())
        }
    }
}

func main() {
    repl()
}
