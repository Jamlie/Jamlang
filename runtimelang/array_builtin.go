package runtimelang

import (
    "fmt"
    "os"
    "slices"
)

func jamlangArrayPush(arr **[]RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: push takes 1 argument")
            os.Exit(0)
        }

        newArray := append(**arr, args[0])
        *arr = &newArray
        return MakeArrayValue(newArray)
    }, "push")
}

func jamlangArrayPop(arr []RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) == 1 {
            if args[0].Type() != Number {
                fmt.Fprintln(os.Stderr, "Error: pop takes a number as an argument")
                os.Exit(0)
            }

            index := int(args[0].(NumberValue).Value)
            if index < 0 || index >= len(arr) {
                fmt.Fprintln(os.Stderr, "Error: pop index out of bounds")
                os.Exit(0)
            }

            arr = append(arr[:index], arr[index+1:]...)
        } else if len(args) == 0 {
            arr = arr[:len(arr)-1]
        }
        
        return MakeArrayValue(arr)
    }, "pop")
}

func jamlangArrayShift(arr []RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        arr = arr[1:]
        return MakeArrayValue(arr)
    }, "shift")
}

func jamlangArrayContains(arr []RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: contains takes 1 argument")
            os.Exit(0)
        }

        return MakeBoolValue(slices.Contains(arr, args[0]))
    }, "contains")
}

func jamlangArrayInsertInto(arr []RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: insert takes 2 arguments")
            os.Exit(0)
        }

        if args[0].Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: insert takes a number as an argument")
            os.Exit(0)
        }

        index := int(args[0].(NumberValue).Value)
        if index < 0 || index >= len(arr) {
            fmt.Fprintln(os.Stderr, "Error: insert index out of bounds")
            os.Exit(0)
        }

        arr = append(arr[:index], append([]RuntimeValue{args[1]}, arr[index:]...)...)
        return MakeArrayValue(arr)
    }, "insertInto")
}

func jamlangArrayPushAll(arr []RuntimeValue) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: pushAll takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != Array {
            fmt.Fprintln(os.Stderr, "Error: pushAll takes an array as an argument")
            os.Exit(0)
        }

        arr = append(arr, args[0].(ArrayValue).Values...)
        return MakeArrayValue(arr)
    }, "pushAll")
}
