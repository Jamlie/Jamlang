package runtimelang

import (
    "fmt"
    "os"
)

type Environment struct {
    parent    *Environment
    variables map[string]RuntimeValue
    constants map[string]bool
}

func CreateGlobalEnvironment() *Environment {
    env := NewEnvironment(nil)
    env.DeclareVariable("true", MakeBoolValue(true), true)
    env.DeclareVariable("false", MakeBoolValue(false), true)
    env.DeclareVariable("null", MakeNullValue(), true)

    timeObject := make(map[string]RuntimeValue)
    timeObject["now"] = MakeNativeFunction(jamlangCurrentTime, "now")
    timeObject["sleep"] = MakeNativeFunction(jamlangSleep, "sleep")
    env.DeclareVariable("Time", MakeObjectValue(timeObject), true)

    bitwiseObject := make(map[string]RuntimeValue)
    bitwiseObject["NOT"] = MakeNativeFunction(jamlangBitwiseNot, "NOT")
    bitwiseObject["AND"] = MakeNativeFunction(jamlangBitwiseAnd, "AND")
    bitwiseObject["OR"] = MakeNativeFunction(jamlangBitwiseOr, "OR")
    bitwiseObject["XOR"] = MakeNativeFunction(jamlangBitwiseXor, "XOR")
    env.DeclareVariable("Bitwise", MakeObjectValue(bitwiseObject), true)

    // osObject := make(map[string]RuntimeValue)
    // osObject["exit"] = MakeNativeFunction(jamlangExit)
    // osObject["open"] = MakeNativeFunction(jamlangOpen)
    // osObject["read"] = MakeNativeFunction(jamlangRead)
    // osObject["write"] = MakeNativeFunction(jamlangWrite)
    // osObject["close"] = MakeNativeFunction(jamlangClose)
    // osObject["remove"] = MakeNativeFunction(jamlangRemove)
    // osObject["rename"] = MakeNativeFunction(jamlangRename)
    // env.DeclareVariable("OS", MakeObjectValue(osObject), true)

    env.DeclareVariable("println", MakeNativeFunction(jamlangPrintln, "println"), true)
    env.DeclareVariable("print", MakeNativeFunction(jamlangPrint, "print"), true)
    env.DeclareVariable("typeof", MakeNativeFunction(jamlangTypeof, "typeof"), true)
    env.DeclareVariable("exit", MakeNativeFunction(jamlangExit, "exit"), true)
    env.DeclareVariable("input", MakeNativeFunction(jamlangInput, "input"), true)
    env.DeclareVariable("len", MakeNativeFunction(jamlangLen, "len"), true)
    env.DeclareVariable("append", MakeNativeFunction(jamlangAppend, "append"), true)
    env.DeclareVariable("pop", MakeNativeFunction(jamlangPop, "pop"), true)
    env.DeclareVariable("array", MakeNativeFunction(jamlangArray, "array"), true)
    env.DeclareVariable("tuple", MakeNativeFunction(jamlangTuple, "tuple"), true)
    env.DeclareVariable("hex", MakeNativeFunction(jamlangHex, "hex"), true)
    env.DeclareVariable("string", MakeNativeFunction(jamlangToString, "string"), true)
    env.DeclareVariable("uint32", MakeNativeFunction(jamlangToUint32, "int"), true)
    env.DeclareVariable("uint64", MakeNativeFunction(jamlangToUint64, "int"), true)
    env.DeclareVariable("int32", MakeNativeFunction(jamlangToInt32, "int"), true)
    env.DeclareVariable("int64", MakeNativeFunction(jamlangToInt64, "int"), true)
    env.DeclareVariable("float", MakeNativeFunction(jamlangToFloat, "float"), true)
    env.DeclareVariable("eval", MakeNativeFunction(jamlangEval, "eval"), true)

    return env
}

func NewEnvironment(parent *Environment) *Environment {
    return &Environment{
        parent:    parent,
        variables: make(map[string]RuntimeValue),
        constants: make(map[string]bool),
    }
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, constant bool) RuntimeValue {
    if _, ok := e.variables[name]; ok {
        if _, ok := e.variables[name].(FunctionValue); ok {
            fmt.Printf("Function %s already declared\n", name)
        } else {
            fmt.Printf("Variable %s already declared\n", name)
        }
        os.Exit(0)
        return nil
    }

    e.variables[name] = value

    if constant {
        e.constants[name] = true
    }

    return value
}

func (e *Environment) AssignVariable(name string, value RuntimeValue) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Printf("Variable %s not declared\n", name)
        os.Exit(0)
        return nil
    }

    if env.constants[name] {
        fmt.Printf("Variable %s is constant. Cannot reassign a constant.\n", name)
        os.Exit(0)
        return nil
    }

    env.variables[name] = value

    return value
}

func (e *Environment) Resolve(name string) *Environment {
    if _, ok := e.variables[name]; ok {
        return e
    }

    if e.parent == nil {
        return nil
    }

    return e.parent.Resolve(name)
}

func (e *Environment) LookupVariable(name string) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Printf("%s not declared\n", name)
        os.Exit(0)
        return nil
    }

    return env.variables[name]
}

func (e *Environment) RemoveVariable(name string) {
    delete(e.variables, name)
}
