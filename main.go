package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Jamlee977/CustomLanguage/parser"
	"github.com/Jamlee977/CustomLanguage/runtimelang"
)

func repl() {
    parser := parser.NewParser()
    file, err := os.Open("test.txt")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        text := scanner.Text()
        if text == "" {
            continue
        }
        program := parser.ProduceAST(text)
        runtimeValue, err := runtimelang.Evaluate(&program)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        fmt.Println(runtimeValue.Get())
    }
}

func main() {
    repl()
}
