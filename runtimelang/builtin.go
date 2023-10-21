package runtimelang

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "time"
    // "strings"

    "github.com/Jamlie/Jamlang/parser"
)

func jamlangPrintln(args []RuntimeValue, environment Environment) RuntimeValue {
    for _, arg := range args {
        fmt.Print(arg.Get())
    }
    fmt.Println()
    return MakeNullValue()
}

func jamlangPrint(args []RuntimeValue, environment Environment) RuntimeValue {
    for _, arg := range args {
        fmt.Print(arg.Get())
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

    fmt.Print(args[0].Get())
    scanner := bufio.NewReader(os.Stdin)
    input, err := scanner.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input")
        os.Exit(0)
    }
    return MakeStringValue(input)
}

func jamlangLen(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("len takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == "array" {
        goArray := ToGoArrayValue(args[0].(ArrayValue))
        return MakeNumberValue(float64(len(goArray)))
    } else if args[0].Type() == "tuple" {
        goTuple := ToGoTupleValue(args[0].(TupleValue))
        return MakeNumberValue(float64(len(goTuple)))
    } else if args[0].Type() == "string" {
        goString := ToGoStringValue(args[0].(StringValue))
        return MakeNumberValue(float64(len(goString)))
    } else {
        fmt.Println("len takes an array, tuple or string")
        os.Exit(0)
        return MakeNullValue()
    }
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

func jamlangSetArrayElement(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 3 {
        fmt.Println("set takes 3 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "array" {
        fmt.Println("set takes an array")
        os.Exit(0)
    }

    goArray := ToGoArrayValue(args[0].(ArrayValue))
    index := int(ToGoNumberValue(args[1].(NumberValue)))
    if index < 0 || index >= len(goArray) {
        fmt.Println("set takes a valid index")
        os.Exit(0)
    }

    goArray[index] = args[2]
    return MakeArrayValue(goArray)
}

func jamlangCopy(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("copy takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == "array" {
        goArray := ToGoArrayValue(args[0].(ArrayValue))
        goArrayCopy := make([]RuntimeValue, len(goArray))
        copy(goArrayCopy, goArray)
        return MakeArrayValue(goArrayCopy)
    } else if args[0].Type() == "tuple" {
        goTuple := ToGoTupleValue(args[0].(TupleValue))
        goTupleCopy := make([]RuntimeValue, len(goTuple))
        copy(goTupleCopy, goTuple)
        return MakeTupleValue(goTupleCopy)
    } else {
        fmt.Println("copy takes an array or tuple")
        os.Exit(0)
        return MakeNullValue()
    }
}

func jamlangTuple(args []RuntimeValue, environment Environment) RuntimeValue {
    return MakeTupleValue(args)
}

func jamlangToString(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("string takes 1 argument")
        os.Exit(0)
    }

    return MakeStringValue(args[0].ToString())
}

func jamlangHex(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("hex takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "string" {
        fmt.Println("hex takes a string")
        os.Exit(0)
    }

    hexString := args[0].ToString()

    hexInt, err := strconv.ParseInt(hexString, 16, 64)
    if err != nil {
        fmt.Println("hex takes a string")
        os.Exit(0)
    }

    return MakeNumberValue(float64(hexInt))
}

func jamlangCurrentTime(args []RuntimeValue, environment Environment) RuntimeValue {
    return MakeNumberValue(float64(time.Now().UnixMicro()))
}

func jamlangBitwiseNot(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("bitwise not takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "number" {
        fmt.Println("bitwise not takes a number")
        os.Exit(0)
    }

    return MakeNumberValue(float64(^int64(args[0].(NumberValue).Value)))
}

func jamlangBitwiseAnd(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("bitwise and takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "number" || args[1].Type() != "number" {
        fmt.Println("bitwise and takes 2 numbers")
        os.Exit(0)
    }

    return MakeNumberValue(float64(int64(args[0].(NumberValue).Value) & int64(args[1].(NumberValue).Value)))
}

func jamlangBitwiseOr(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("bitwise or takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "number" || args[1].Type() != "number" {
        fmt.Println("bitwise or takes 2 numbers")
        os.Exit(0)
    }

    return MakeNumberValue(float64(int64(args[0].(NumberValue).Value) | int64(args[1].(NumberValue).Value)))
}

func jamlangBitwiseXor(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("bitwise xor takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "number" || args[1].Type() != "number" {
        fmt.Println("bitwise xor takes 2 numbers")
        os.Exit(0)
    }

    return MakeNumberValue(float64(int64(args[0].(NumberValue).Value) ^ int64(args[1].(NumberValue).Value)))
}

func jamlangEval(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("eval takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "string" {
        fmt.Println("eval takes a string")
        os.Exit(0)
    }

    code := args[0].ToString()
    program := parser.NewParser().ProduceAST(code)
    newEnvironment := CreateGlobalEnvironment()
    Evaluate(&program, *newEnvironment)
    return MakeNullValue()
}

func jamlangOpen(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("open takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "string" {
        fmt.Println("open takes a string")
        os.Exit(0)
    }

    filename := args[0].(StringValue).Value

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening file")
        os.Exit(0)
    }

    return MakeFileValue(file.Name(), file)
}

func jamlangClose(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Println("close takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "file" {
        fmt.Println("close takes a file")
        os.Exit(0)
    }

    if args[0].(FileValue).IsClosed {
        fmt.Println("File is already closed")
        os.Exit(0)
    }

    f, _ := args[0].(FileValue)
    f.IsClosed = true
    args[0].(FileValue).File.Close()
    f.File.Close()
    return MakeNullValue()
}

func jamlangWrite(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Println("write takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "file" {
        fmt.Println("write takes a file")
        os.Exit(0)
    }

    if args[0].(FileValue).IsClosed {
        fmt.Println("File is closed")
        os.Exit(0)
    }

    if args[1].Type() != "string" {
        fmt.Println("write takes a string")
        os.Exit(0)
    }

    file := ToGoFileValue(args[0].(FileValue))

    file.Truncate(0)
    file.Seek(0, 0)

    file.WriteString(args[1].(StringValue).ToString())

    return MakeNullValue()
}
