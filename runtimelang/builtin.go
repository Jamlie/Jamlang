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

func jamlangArray(args []RuntimeValue, environment Environment) RuntimeValue {
    return MakeArrayValue(args)
}

func jamlangLen(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("len takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "array" {
        fmt.Println("len takes an array")
        os.Exit(0)
    }

    goArray := ToGoArrayValue(args[0].(ArrayValue))
    return MakeNumberValue(float64(len(goArray)))
}

func jamlangAppend(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("append takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "array" {
        fmt.Println("append takes an array")
        os.Exit(0)
    }

    goArray := ToGoArrayValue(args[0].(ArrayValue))
    goArray = append(goArray, args[1])
    return MakeArrayValue(goArray)
}

func jamlangPop(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("pop takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "array" {
        fmt.Println("pop takes an array")
        os.Exit(0)
    }

    goArray := ToGoArrayValue(args[0].(ArrayValue))
    if len(goArray) == 0 {
        fmt.Println("pop takes a non-empty array")
        os.Exit(0)
    }

    return MakeArrayValue(goArray[:len(goArray)-1])
}

func jamlangGetArrayElement(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("get takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "array" {
        fmt.Println("get takes an array")
        os.Exit(0)
    }

    goArray := ToGoArrayValue(args[0].(ArrayValue))
    index := int(ToGoNumberValue(args[1].(NumberValue)))
    if index < 0 || index >= len(goArray) {
        fmt.Println("get takes a valid index")
        os.Exit(0)
    }

    return goArray[index]
}
