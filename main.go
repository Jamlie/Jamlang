package main

import (
    "fmt"
    "os"
    "flag"
    "strings"
    "bufio"
    "io/ioutil"
    
    "github.com/Jamlee977/CustomLanguage/parser"
    "github.com/Jamlee977/CustomLanguage/runtimelang"
)

var (
    env = runtimelang.CreateGlobalEnvironment()
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
            os.Exit(0)
        }
        fmt.Println(runtimeValue.Get())
    }
}

func main() {
    runFlag := flag.Bool("r", false, "Run a file")
    flag.Parse()
    args := flag.Args()
    if flag.NFlag() == 0 && len(args) == 0 {
        repl()
    } else {
        if *runFlag {
            if len(args) < 1 {
                fmt.Println("No file specified")
                os.Exit(0)
            }

            if !strings.HasSuffix(args[0], ".jam") {
                fmt.Println("File must have .jam extension")
                os.Exit(0)
            }

            data, err := ioutil.ReadFile(args[0])
            if err != nil {
                fmt.Println(err)
                os.Exit(0)
            }

            program := parser.NewParser().ProduceAST(string(data))
            _, err = runtimelang.Evaluate(&program, *env)
            if err != nil {
                fmt.Println(err)
                os.Exit(0)
            }
        } else {
             option := args[0]
            if option == "run" {
                if len(args) < 2 {
                    fmt.Println("No file specified")
                    os.Exit(0)
                }

                if !strings.HasSuffix(args[1], ".jam") {
                    fmt.Println("File must have .jam extension")
                    os.Exit(0)
                }

                data, err := ioutil.ReadFile(args[1])
                if err != nil {
                    fmt.Println(err)
                    os.Exit(0)
                }

                parser := parser.NewParser()
                program := parser.ProduceAST(string(data))

                _, err = runtimelang.Evaluate(&program, *env)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(0)
                }

            } else {
                fmt.Println("Unknown option")
                fmt.Println("Usage: jamlang [run] [file]")
                fmt.Println("Usage: jamlang -r [file]")
                fmt.Println("Usage: jamlang")
            }
        }
    }
}

/*
{
    class := ClassValue{
        Name: expr.Name,
        Methods: make(map[string]FunctionValue),
    }
    for _, method := range expr.Body {
        if method.Kind() == ast.FunctionDeclarationType {
            method := method.(*ast.FunctionDeclaration)
            if method.Name == "constructor" {
                class.Constructor = FunctionValue{
                    Name: method.Name,
                    Parameters: method.Parameters,
                    Body: method.Body,
                    DeclarationEnvironment: *env,
                }

                continue
            }
            class.Methods[method.Name] = FunctionValue{
                Name: method.Name,
                Parameters: method.Parameters,
                Body: method.Body,
                DeclarationEnvironment: *env,
            }

            env.DeclareVariable(method.Name, class.Methods[method.Name], true)
        } else if method.Kind() == ast.VariableDeclarationType {
            method := method.(*ast.VariableDeclaration)
            value, err := Evaluate(method.Value, *env)
            if err != nil {
                return MakeNullValue(), err
            }
            env.DeclareVariable(method.Identifier, value, true)
            continue
        } else {
            return MakeNullValue(), fmt.Errorf("invalid method declaration")
        }
    }
    env.DeclareVariable(expr.Name, class.Constructor, true)
    return MakeNullValue(), nil
}
*/
