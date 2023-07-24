package main

import (
	"bufio"
	"fmt"
	"os"
    "flag"
    "strings"

	"github.com/Jamlee977/CustomLanguage/parser"
	"github.com/Jamlee977/CustomLanguage/runtimelang"
)

var (
    env = runtimelang.NewEnvironment(nil)
)

func repl() {
    parser := parser.NewParser()
    fmt.Println("Repl mode. Type 'exit' to exit.")
    for {
        fmt.Print("> ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        text := scanner.Text()
        if text == "exit" {
            break
        }
        if text == "" {
            continue
        }

        program := parser.ProduceAST(text)
        runtimeValue, err := runtimelang.Evaluate(&program, *env)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        fmt.Println(runtimeValue.Get())
    }
}

func main() {
    env.DeclareVariable("null", runtimelang.MakeNullValue(), true)
    env.DeclareVariable("true", runtimelang.MakeBoolValue(true), true)
    env.DeclareVariable("false", runtimelang.MakeBoolValue(false), true)
    flag.Parse()
    args := flag.Args()
    if len(args) == 0 {
        repl()
    } else {
        option := args[0]
        if option == "run" {
            if len(args) < 2 {
                fmt.Println("No file specified")
                os.Exit(1)
            }

            if !strings.HasSuffix(args[1], ".jam") {
                fmt.Println("File must have .jam extension")
                os.Exit(1)
            }
            
            file, err := os.Open(args[1])
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
                parser := parser.NewParser()
                program := parser.ProduceAST(text)
                runtimeValue, err := runtimelang.Evaluate(&program, *env)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                fmt.Println(runtimeValue.Get())
            }
        } else {
            fmt.Println("Unknown option")
            fmt.Println("Usage: elang [run] [file]")
        }
    }
}
