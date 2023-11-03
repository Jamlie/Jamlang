package runtimelang

import (
    "fmt"
    "os"
    "strings"
)

func jamlangStringToUpper(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 0 {
            fmt.Fprintln(os.Stderr, "Error: toUpper takes 0 arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.ToUpper(str)) 
    }, "toUpper")
}

func jamlangStringToLower(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 0 {
            fmt.Fprintln(os.Stderr, "Error: toLower takes 0 arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.ToLower(str)) 
    }, "toLower")
}

func jamlangStringContains(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: contains takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: contains takes a string as an argument")
            os.Exit(0)
        }

        return MakeBoolValue(strings.Contains(str, args[0].(StringValue).Value)) 
    }, "contains")
}

func jamlangStringSplit(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: split takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: split takes a string as an argument")
            os.Exit(0)
        }

        var split []RuntimeValue
        for _, s := range strings.Split(str, args[0].(StringValue).Value) {
            split = append(split, MakeStringValue(s))
        }
        return MakeArrayValue(split)
    }, "split")
}

func jamlangStringEqualsIgnoreCase(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: equalsIgnoreCase takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: equalsIgnoreCase takes a string as an argument")
            os.Exit(0)
        }

        return MakeBoolValue(strings.EqualFold(str, args[0].(StringValue).Value)) 
    }, "equalsIgnoreCase")
}

func jamlangStringStartsWith(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: startsWith takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: startsWith takes a string as an argument")
            os.Exit(0)
        }

        return MakeBoolValue(strings.HasPrefix(str, args[0].(StringValue).Value)) 
    }, "startsWith")
}

func jamlangStringEndsWith(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: endsWith takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: endsWith takes a string as an argument")
            os.Exit(0)
        }

        return MakeBoolValue(strings.HasSuffix(str, args[0].(StringValue).Value)) 
    }, "endsWith")
}

func jamlangStringIndexOf(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: indexOf takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: indexOf takes a string as an argument")
            os.Exit(0)
        }

        return MakeNumberValue(float64(strings.Index(str, args[0].(StringValue).Value)))
    }, "indexOf")
}

func jamlangStringLastIndexOf(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: lastIndexOf takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: lastIndexOf takes a string as an argument")
            os.Exit(0)
        }

        return MakeNumberValue(float64(strings.LastIndex(str, args[0].(StringValue).Value)))
    }, "lastIndexOf")
}

func jamlangStringSubstring(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: substring takes 2 arguments")
            os.Exit(0)
        }

        if args[0].Type() != Number || args[1].Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: substring takes 2 numbers as arguments")
            os.Exit(0)
        }

        return MakeStringValue(str[int(args[0].(NumberValue).Value):int(args[1].(NumberValue).Value)])
    }, "substring")
}

func jamlangStringReplace(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: replace takes 2 arguments")
            os.Exit(0)
        }

        if args[0].Type() == Number {
            return MakeStringValue(strings.Replace(str, string(rune(args[0].(NumberValue).Value)), args[1].(StringValue).Value, -1))
        }

        if args[0].Type() != String || args[1].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: replace takes 2 strings or a number and a string as arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.Replace(str, args[0].(StringValue).Value, args[1].(StringValue).Value, -1))
    }, "replace")
}

func jamlangStringTrim(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 0 {
            fmt.Fprintln(os.Stderr, "Error: trim takes 0 arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.TrimSpace(str))
    }, "trim")
}

func jamlangStringTrimLeft(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 0 {
            fmt.Fprintln(os.Stderr, "Error: trimLeft takes 0 arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.TrimLeftFunc(str, func(r rune) bool {
            return r == ' ' || r == '\t' || r == '\n' || r == '\r'
        }))
    }, "trimLeft")
}

func jamlangStringTrimRight(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 0 {
            fmt.Fprintln(os.Stderr, "Error: trimRight takes 0 arguments")
            os.Exit(0)
        }

        return MakeStringValue(strings.TrimRightFunc(str, func(r rune) bool {
            return r == ' ' || r == '\t' || r == '\n' || r == '\r'
        }))
    }, "trimRight")
}

func jamlangStringRepeat(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "Error: repeat takes 1 argument")
            os.Exit(0)
        }

        if args[0].Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: repeat takes a number as an argument")
            os.Exit(0)
        }

        return MakeStringValue(strings.Repeat(str, int(args[0].(NumberValue).Value)))
    }, "repeat")
}

func jamlangStringLeftPad(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: leftPad takes 2 arguments")
            os.Exit(0)
        }

        if args[0].Type() != Number || args[1].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: leftPad takes a number and a string as arguments")
            os.Exit(0)
        }

        return MakeStringValue(fmt.Sprintf("%"+args[0].(NumberValue).ToString()+"s", str))
    }, "leftPad")
}

func jamlangStringRightPad(str string) RuntimeValue {
    return MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: rightPad takes 2 arguments")
            os.Exit(0)
        }

        if args[0].Type() != Number || args[1].Type() != String {
            fmt.Fprintln(os.Stderr, "Error: rightPad takes a number and a string as arguments")
            os.Exit(0)
        }

        return MakeStringValue(fmt.Sprintf("%-"+args[0].(NumberValue).ToString()+"s", str))
    }, "rightPad")
}

