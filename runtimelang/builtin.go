package runtimelang

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "time"

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
        fmt.Fprintln(os.Stderr, "Error: sleep takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 {
        fmt.Fprintln(os.Stderr, "Error: sleep takes a number - time in milliseconds")
        os.Exit(0)
    }

    time.Sleep(time.Duration(args[0].Get().(float64)) * time.Millisecond)
    return MakeNullValue()
}

func jamlangTypeof(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: typeof takes 1 argument")
        os.Exit(0)
    }

    return MakeStringValue(string(args[0].Type()))
}

func jamlangExit(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: exit takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 {
        fmt.Fprintln(os.Stderr, "Error: exit takes a number - exit code")
        os.Exit(0)
    }

    os.Exit(args[0].(IntValue).GetInt())
    return MakeNullValue()
}

func jamlangInput(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: input takes 1 argument")
        os.Exit(0)
    }

    fmt.Print(args[0].Get())
    scanner := bufio.NewReader(os.Stdin)
    input, err := scanner.ReadString('\n')
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: reading input")
        os.Exit(0)
    }
    return MakeStringValue(input)
}

// func jamlangLen(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 1 {
//         fmt.Fprintln(os.Stderr, "len takes 1 argument")
//         os.Exit(0)
//     }
//
//     if args[0].Type() == Array {
//         goArray := ToGoArrayValue(args[0].(ArrayValue))
//         return MakeInt64Value(int64(len(goArray)))
//     } else if args[0].Type() == Tuple {
//         goTuple := ToGoTupleValue(args[0].(TupleValue))
//         return MakeInt(float64(len(goTuple)))
//     } else if args[0].Type() == String {
//         goString := ToGoStringValue(args[0].(StringValue))
//         return MakeNumberValue(float64(len(goString)))
//     } else {
//         fmt.Println("len takes an array, tuple or string")
//         os.Exit(0)
//         return MakeNullValue()
//     }
// }

// func jamlangAppend(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 2 {
//         fmt.Fprintln(os.Stderr, "append takes 2 arguments")
//         os.Exit(0)
//     }
//
//     if args[0].Type() != Array {
//         fmt.Fprintln(os.Stderr, "append takes an array")
//         os.Exit(0)
//     }
//
//     goArray := ToGoArrayValue(args[0].(ArrayValue))
//     goArray = append(goArray, args[1])
//     return MakeArrayValue(goArray)
// }

// func jamlangPop(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 1 {
//         fmt.Fprintln(os.Stderr, "pop takes 1 argument")
//         os.Exit(0)
//     }
//
//     if args[0].Type() != Array {
//         fmt.Fprintln(os.Stderr, "pop takes an array")
//         os.Exit(0)
//     }
//
//     goArray := ToGoArrayValue(args[0].(ArrayValue))
//     if len(goArray) == 0 {
//         fmt.Fprintln(os.Stderr, "pop takes a non-empty array")
//         os.Exit(0)
//     }
//
//     return MakeArrayValue(goArray[:len(goArray)-1])
// }

func jamlangCopy(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: copy takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == Array {
        goArray := ToGoArrayValue(args[0].(ArrayValue))
        goArrayCopy := make([]RuntimeValue, len(goArray))
        copy(goArrayCopy, goArray)
        return MakeArrayValue(goArrayCopy)
    } else if args[0].Type() == Tuple {
        goTuple := ToGoTupleValue(args[0].(TupleValue))
        goTupleCopy := make([]RuntimeValue, len(goTuple))
        copy(goTupleCopy, goTuple)
        return MakeTupleValue(goTupleCopy)
    } else {
        fmt.Fprintln(os.Stderr, "Error: copy takes an array or tuple")
        os.Exit(0)
        return MakeNullValue()
    }
}

func jamlangArray(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: array takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 {
        fmt.Fprintln(os.Stderr, "Error: array takes a number")
        os.Exit(0)
    }

    var size int
    switch args[0].Type() {
    case I8:
        size = int(ToGoNumberValue(args[0].(Int8Value)))
    case I16:
        size = int(ToGoNumberValue(args[0].(Int16Value)))
    case I32:
        size = int(ToGoNumberValue(args[0].(Int32Value)))
    default:
        fmt.Fprintln(os.Stderr, "Error: array takes a number")
        os.Exit(0)
    }

    if size < 0 {
        fmt.Fprintln(os.Stderr, "Error: array takes a positive number")
        os.Exit(0)
    }

    goArray := make([]RuntimeValue, size)
    return MakeArrayValue(goArray)
}

func jamlangTuple(args []RuntimeValue, environment Environment) RuntimeValue {
    return MakeTupleValue(args)
}

func jamlangToString(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: string takes 1 argument")
        os.Exit(0)
    }

    return MakeStringValue(args[0].ToString())
}

func jamlangToUint32(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: uint32 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int8Value)))))
    }

    if args[0].Type() == I16 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int16Value)))))
    }

    if args[0].Type() == I32 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int32Value)))))
    }

    if args[0].Type() == I64 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int64Value)))))
    }

    if args[0].Type() == F32 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Float32Value)))))
    }

    if args[0].Type() == F64 {
        return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Float64Value)))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: uint32 takes a string or a number")
        os.Exit(0)
    }

    uintString := args[0].ToString()

    uintUint, err := strconv.ParseUint(uintString, 10, 32)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: uint32 takes a string or a number")
        os.Exit(0)
    }

    return MakeInt64Value(int64(uintUint))
}

func jamlangToUint64(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: uint64 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int8Value)))))
    }

    if args[0].Type() == I16 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int16Value)))))
    }

    if args[0].Type() == I32 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int32Value)))))
    }

    if args[0].Type() == I64 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int64Value)))))
    }

    if args[0].Type() == F32 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Float32Value)))))
    }

    if args[0].Type() == F64 {
        return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Float64Value)))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: uint64 takes a string or a number")
        os.Exit(0)
    }

    uintString := args[0].ToString()

    uintUint, err := strconv.ParseUint(uintString, 10, 64)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: uint64 takes a string or a number")
        os.Exit(0)
    }

    return MakeFloat32Value(float32(uintUint))
}

func jamlangToInt8(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: int8 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeInt8Value(int8(ToGoNumberValue(args[0].(Float64Value))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: int8 takes a string or a number")
        os.Exit(0)
    }

    intString := args[0].ToString()

    intInt, err := strconv.ParseInt(intString, 10, 8)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: int8 takes a string or a number")
        os.Exit(0)
    }

    return MakeInt8Value(int8(intInt))
}

func jamlangToInt16(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: int16 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeInt16Value(int16(ToGoNumberValue(args[0].(Float64Value))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: int16 takes a string or a number")
        os.Exit(0)
    }

    intString := args[0].ToString()

    intInt, err := strconv.ParseInt(intString, 10, 16)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: int16 takes a string or a number")
        os.Exit(0)
    }

    return MakeInt16Value(int16(intInt))
}

func jamlangToInt32(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: int32 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeInt32Value(int32(ToGoNumberValue(args[0].(Float64Value))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: int32 takes a string or a number")
        os.Exit(0)
    }

    intString := args[0].ToString()

    intInt, err := strconv.ParseInt(intString, 10, 32)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: int32 takes a string or a number")
        os.Exit(0)
    }

    return MakeInt32Value(int32(intInt))
}

func jamlangToInt64(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: int64 takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeInt64Value(int64(ToGoNumberValue(args[0].(Float64Value))))
    }


    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: int64 takes a string or a number")
        os.Exit(0)
    }

    intString := args[0].ToString()

    intInt, err := strconv.ParseInt(intString, 10, 64)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: int64 takes a string or a number")
        os.Exit(0)
    }

    return MakeInt64Value(int64(intInt))
}

func jamlangToFloat32(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Float64Value))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
        os.Exit(0)
    }

    floatString := args[0].ToString()

    floatFloat, err := strconv.ParseFloat(floatString, 32)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
        os.Exit(0)
    }

    return MakeFloat32Value(float32(floatFloat))
}

func jamlangToFloat64(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
        os.Exit(0)
    }

    if args[0].Type() == I8 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int8Value))))
    }

    if args[0].Type() == I16 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int16Value))))
    }

    if args[0].Type() == I32 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int32Value))))
    }

    if args[0].Type() == I64 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int64Value))))
    }

    if args[0].Type() == F32 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Float32Value))))
    }

    if args[0].Type() == F64 {
        return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Float64Value))))
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
        os.Exit(0)
    }

    floatString := args[0].ToString()

    floatFloat, err := strconv.ParseFloat(floatString, 64)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
        os.Exit(0)
    }

    return MakeFloat64Value(float64(floatFloat))
}

func jamlangHex(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: hex takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "Error: hex takes a string")
        os.Exit(0)
    }

    hexString := args[0].ToString()

    hexInt, err := strconv.ParseInt(hexString, 16, 64)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: hex takes a string")
        os.Exit(0)
    }

    return MakeInt64Value(int64(hexInt))
}

func jamlangCurrentTime(args []RuntimeValue, environment Environment) RuntimeValue {
    return MakeInt64Value(time.Now().UnixMicro())
}

func jamlangBitwiseNot(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64 {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes a number")
        os.Exit(0)
    }

    switch args[0].Type() {
    case I8:
        return MakeInt8Value(^int8(ToGoNumberValue(args[0].(Int8Value))))
    case I16:
        return MakeInt16Value(^int16(ToGoNumberValue(args[0].(Int16Value))))
    case I32:
        return MakeInt32Value(^int32(ToGoNumberValue(args[0].(Int32Value))))
    case I64:
        return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Int64Value))))
    case F32:
        return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Float32Value))))
    case F64:
        return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Float64Value))))
    default:
        fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes a number")
        os.Exit(0)
        return nil
    }
}

func jamlangBitwiseAnd(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 arguments")
        os.Exit(0)
    }

    if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.AND and takes 2 numbers")
        os.Exit(0)
    }

    switch args[0].Type() {
    case I8:
        switch args[1].Type() {
        case I8:
            return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) & ToGoNumberValue(args[1].(Int8Value))))
        case I16:
            return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) & int8(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) & int8(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I16:
        switch args[1].Type() {
        case I8:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) & ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) & ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) & ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I32:
        switch args[1].Type() {
        case I8:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) & ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F32:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) & int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) & int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) & int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) & int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) & int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) & int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) & int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) & int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) & int64(ToGoNumberValue(args[1].(Float64Value))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
            os.Exit(0)
            return nil
        }
    default:
        fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
        os.Exit(0)
        return nil
    }
}

func jamlangBitwiseOr(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 arguments")
        os.Exit(0)
    }

    if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
        os.Exit(0)
    }

    switch args[0].Type() {
    case I8:
        switch args[1].Type() {
        case I8:
            return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) | ToGoNumberValue(args[1].(Int8Value))))
        case I16:
            return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) | int8(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) | int8(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I16:
        switch args[1].Type() {
        case I8:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) | ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) | ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) | ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I32:
        switch args[1].Type() {
        case I8:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) | ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F32:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) | int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) | int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) | int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) | int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) | int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) | int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) | int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) | int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) | int64(ToGoNumberValue(args[1].(Float64Value))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    default:
        fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
        os.Exit(0)
        return nil
    }
}

func jamlangBitwiseXor(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 2 {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 arguments")
        os.Exit(0)
    }

    if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
        fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
        os.Exit(0)
    }

    switch args[0].Type() {
    case I8:
        switch args[1].Type() {
        case I8:
            return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) ^ ToGoNumberValue(args[1].(Int8Value))))
        case I16:
            return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) ^ int8(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) ^ int8(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I16:
        switch args[1].Type() {
        case I8:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) ^ ToGoNumberValue(args[1].(Int16Value))))
        case I32:
            return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) ^ ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I32:
        switch args[1].Type() {
        case I8:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ ToGoNumberValue(args[1].(Int32Value))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case I64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ ToGoNumberValue(args[1].(Int64Value))))
        case F32:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F32:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) ^ int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) ^ int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) ^ int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Float64Value)))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    case F64:
        switch args[1].Type() {
        case I8:
            return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) ^ int8(ToGoNumberValue(args[1].(Int8Value)))))
        case I16:
            return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) ^ int16(ToGoNumberValue(args[1].(Int16Value)))))
        case I32:
            return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) ^ int32(ToGoNumberValue(args[1].(Int32Value)))))
        case I64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) ^ int64(ToGoNumberValue(args[1].(Int64Value)))))
        case F32:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
        case F64:
            return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) ^ int64(ToGoNumberValue(args[1].(Float64Value))))
        default:
            fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
            os.Exit(0)
            return nil
        }
    default:
        fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
        os.Exit(0)
        return nil
    }
}

func jamlangEval(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "Error: eval takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "eval takes a string")
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
        fmt.Fprintln(os.Stderr, "open takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != String {
        fmt.Fprintln(os.Stderr, "open takes a string")
        os.Exit(0)
    }

    filename := args[0].(StringValue).Value

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error opening file")
        os.Exit(0)
    }

    return MakeFileValue(file.Name(), file)
}

func jamlangClose(args []RuntimeValue, environment Environment) RuntimeValue {
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "close takes 1 argument")
        os.Exit(0)
    }

    if args[0].Type() != "file" {
        fmt.Fprintln(os.Stderr, "close takes a file")
        os.Exit(0)
    }

    if args[0].(FileValue).IsClosed {
        fmt.Fprintln(os.Stderr, "File is already closed")
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
        fmt.Fprintln(os.Stderr, "write takes 2 arguments")
        os.Exit(0)
    }

    if args[0].Type() != "file" {
        fmt.Fprintln(os.Stderr, "write takes a file")
        os.Exit(0)
    }

    if args[0].(FileValue).IsClosed {
        fmt.Fprintln(os.Stderr, "File is closed")
        os.Exit(0)
    }

    if args[1].Type() != String {
        fmt.Fprintln(os.Stderr, "write takes a string")
        os.Exit(0)
    }

    file := ToGoFileValue(args[0].(FileValue))

    file.Truncate(0)
    file.Seek(0, 0)

    file.WriteString(args[1].(StringValue).ToString())

    return MakeNullValue()
}
