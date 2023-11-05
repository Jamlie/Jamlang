package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlie/Jamlang/ast"
)

type Environment struct {
    parent    *Environment
    variables map[string]RuntimeValue
    constants map[string]bool
    types     map[string]ast.VariableType
}

func CreateGlobalEnvironment() *Environment {
    env := NewEnvironment(nil)
    env.DeclareVariable("true", MakeBoolValue(true), true, ast.BoolType)
    env.DeclareVariable("false", MakeBoolValue(false), true, ast.BoolType)
    env.DeclareVariable("null", MakeNullValue(), true, ast.NullType)

    timeObject := make(map[string]RuntimeValue)
    timeObject["now"] = MakeNativeFunction(jamlangCurrentTime, "now")
    timeObject["sleep"] = MakeNativeFunction(jamlangSleep, "sleep")
    env.DeclareVariable("Time", MakeObjectValue(timeObject), true, ast.ObjectType)

    bitwiseObject := make(map[string]RuntimeValue)
    bitwiseObject["NOT"] = MakeNativeFunction(jamlangBitwiseNot, "NOT")
    bitwiseObject["AND"] = MakeNativeFunction(jamlangBitwiseAnd, "AND")
    bitwiseObject["OR"] = MakeNativeFunction(jamlangBitwiseOr, "OR")
    bitwiseObject["XOR"] = MakeNativeFunction(jamlangBitwiseXor, "XOR")
    env.DeclareVariable("Bitwise", MakeObjectValue(bitwiseObject), true, ast.ObjectType)

    // osObject := make(map[string]RuntimeValue)
    // osObject["exit"] = MakeNativeFunction(jamlangExit)
    // osObject["open"] = MakeNativeFunction(jamlangOpen)
    // osObject["read"] = MakeNativeFunction(jamlangRead)
    // osObject["write"] = MakeNativeFunction(jamlangWrite)
    // osObject["close"] = MakeNativeFunction(jamlangClose)
    // osObject["remove"] = MakeNativeFunction(jamlangRemove)
    // osObject["rename"] = MakeNativeFunction(jamlangRename)
    // env.DeclareVariable("OS", MakeObjectValue(osObject), true)

    env.DeclareVariable("println", MakeNativeFunction(jamlangPrintln, "println"), true, ast.FunctionType)
    env.DeclareVariable("print", MakeNativeFunction(jamlangPrint, "print"), true, ast.FunctionType)
    env.DeclareVariable("typeof", MakeNativeFunction(jamlangTypeof, "typeof"), true, ast.StringType)
    env.DeclareVariable("exit", MakeNativeFunction(jamlangExit, "exit"), true, ast.FunctionType)
    env.DeclareVariable("input", MakeNativeFunction(jamlangInput, "input"), true, ast.FunctionType)
    env.DeclareVariable("array", MakeNativeFunction(jamlangArray, "array"), true, ast.ArrayType)
    env.DeclareVariable("tuple", MakeNativeFunction(jamlangTuple, "tuple"), true, ast.TupleType)
    env.DeclareVariable("hex", MakeNativeFunction(jamlangHex, "hex"), true, ast.NumberType)
    env.DeclareVariable("string", MakeNativeFunction(jamlangToString, "string"), true, ast.FunctionType)
    env.DeclareVariable("uint32", MakeNativeFunction(jamlangToUint32, "int"), true, ast.NumberType)
    env.DeclareVariable("uint64", MakeNativeFunction(jamlangToUint64, "int"), true, ast.NumberType)
    env.DeclareVariable("int32", MakeNativeFunction(jamlangToInt32, "int"), true, ast.NumberType)
    env.DeclareVariable("int64", MakeNativeFunction(jamlangToInt64, "int"), true, ast.NumberType)
    env.DeclareVariable("float", MakeNativeFunction(jamlangToFloat, "float"), true, ast.NumberType)
    env.DeclareVariable("eval", MakeNativeFunction(jamlangEval, "eval"), true, ast.AnyType)

    return env
}

func NewEnvironment(parent *Environment) *Environment {
    return &Environment{
        parent:    parent,
        variables: make(map[string]RuntimeValue),
        constants: make(map[string]bool),
        types:     make(map[string]ast.VariableType),
    }
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, constant bool, varType ast.VariableType) RuntimeValue {
    if _, ok := e.variables[name]; ok {
        if _, ok := e.variables[name].(FunctionValue); ok {
            fmt.Fprintf(os.Stderr, "Function %s already declared\n", name)
        } else {
            fmt.Fprintf(os.Stderr, "Variable %s already declared\n", name)
        }
        os.Exit(0)
        return nil
    }

    e.variables[name] = value

    if constant {
        e.constants[name] = true
    }

    e.types[name] = varType

    return value
}

func (e *Environment) AssignVariable(name string, value RuntimeValue) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Fprintf(os.Stderr, "Variable %s not declared\n", name)
        os.Exit(0)
        return nil
    }

    if env.constants[name] {
        fmt.Fprintf(os.Stderr, "Variable %s is constant. Cannot reassign a constant.\n", name)
        os.Exit(0)
        return nil
    }

    env.variables[name] = value

    if env.types[name] != value.VarType() && env.types[name] != ast.AnyType {
        fmt.Fprintf(os.Stderr, "Error: Type mismatch, expected %s got %s\n", env.types[name], value.VarType())
        os.Exit(0)
    }

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
        fmt.Fprintf(os.Stderr, "%s not declared\n", name)
        os.Exit(0)
        return nil
    }

    return env.variables[name]
}

func (e *Environment) RemoveVariable(name string) {
    delete(e.variables, name)
}
