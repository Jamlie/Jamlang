package runtimelang

import (
    "fmt"
    "os"
    "time"
)

func jamlangPrintln(args []RuntimeValue, environment Environment) RuntimeValue { 
    for _, arg := range args {
        fmt.Print(arg.Get(), " ")
    }
    fmt.Println()
    return MakeNullValue()
}

func jamlangPrint(args []RuntimeValue, environment Environment) RuntimeValue { 
    for _, arg := range args {
        fmt.Print(arg.Get(), " ")
    }
    return MakeNullValue()
}

func jamlangSleep(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("sleep takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "number" {
        fmt.Println("sleep takes a number - time in milliseconds")
        os.Exit(0)
    }

    time.Sleep(time.Duration(args[0].Get().(float64)) * time.Millisecond)
    return MakeNullValue()
}

func jamlangTypeof(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("typeof takes 1 argument")
        os.Exit(0)
    }

    return MakeStringValue(string(args[0].Type()))
}

func jamlangExit(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("exit takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "number" {
        fmt.Println("exit takes a number - exit code")
        os.Exit(0)
    }

    os.Exit(int(args[0].Get().(float64)))
    return MakeNullValue()
}

func jamlangInput(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("input takes 1 argument")
        os.Exit(0)
    }

    var input string
    fmt.Scanln(&input)
    return MakeStringValue(input)
}
