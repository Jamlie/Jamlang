package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlie/Jamlang/ast"
)

type Environment struct {
    parent           *Environment
    variables        map[string]RuntimeValue
    constants        map[string]bool
    types            map[string]ast.VariableType
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

    objectObject := make(map[string]RuntimeValue)
    objectObject["keys"] = MakeNativeFunction(jamlangObjectKeys, "keys")
    objectObject["values"] = MakeNativeFunction(jamlangObjectValues, "values")
    objectObject["has"] = MakeNativeFunction(jamlangObjectHas, "has")
    env.DeclareVariable("Object", MakeObjectValue(objectObject), true, ast.ObjectType)

    osObject := make(map[string]RuntimeValue)
    osObject["exit"] = MakeNativeFunction(jamlangExit, "exit")
    osObject["open"] = MakeNativeFunction(jamlangOpen, "open")
    env.DeclareVariable("OS", MakeObjectValue(osObject), true, ast.ObjectType)

    env.DeclareVariable("println", MakeNativeFunction(jamlangPrintln, "println"), true, ast.FunctionType)
    env.DeclareVariable("print", MakeNativeFunction(jamlangPrint, "print"), true, ast.FunctionType)
    env.DeclareVariable("typeof", MakeNativeFunction(jamlangTypeof, "typeof"), true, ast.StringType)
    env.DeclareVariable("exit", MakeNativeFunction(jamlangExit, "exit"), true, ast.FunctionType)
    env.DeclareVariable("input", MakeNativeFunction(jamlangInput, "input"), true, ast.StringType)
    env.DeclareVariable("array", MakeNativeFunction(jamlangArray, "array"), true, ast.ArrayType)
    env.DeclareVariable("tuple", MakeNativeFunction(jamlangTuple, "tuple"), true, ast.TupleType)
    env.DeclareVariable("hex", MakeNativeFunction(jamlangHex, "hex"), true, ast.Int64Type)
    env.DeclareVariable("string", MakeNativeFunction(jamlangToString, "string"), true, ast.FunctionType)
    env.DeclareVariable("uint32", MakeNativeFunction(jamlangToUint32, "uint32"), true, ast.Int64Type)
    env.DeclareVariable("uint64", MakeNativeFunction(jamlangToUint64, "uint64"), true, ast.Float32Type)
    env.DeclareVariable("int8", MakeNativeFunction(jamlangToInt8, "int8"), true, ast.Int8Type)
    env.DeclareVariable("int16", MakeNativeFunction(jamlangToInt16, "int16"), true, ast.Int16Type)
    env.DeclareVariable("int32", MakeNativeFunction(jamlangToInt32, "int32"), true, ast.Int32Type)
    env.DeclareVariable("int64", MakeNativeFunction(jamlangToInt64, "int64"), true, ast.Int64Type)
    env.DeclareVariable("float32", MakeNativeFunction(jamlangToFloat32, "float32"), true, ast.Float32Type)
    env.DeclareVariable("float64", MakeNativeFunction(jamlangToFloat64, "float64"), true, ast.Float64Type)
    env.DeclareVariable("eval", MakeNativeFunction(jamlangEval, "eval"), true, ast.AnyType)

    return env
}

func NewEnvironment(parent *Environment) *Environment {
    return &Environment{
        parent:           parent,
        variables:        make(map[string]RuntimeValue),
        constants:        make(map[string]bool),
        types:            make(map[string]ast.VariableType),
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
        fmt.Fprintf(os.Stderr, "Error: Variable %s not declared\n", name)
        os.Exit(0)
        return nil
    }

    if env.constants[name] {
        fmt.Fprintf(os.Stderr, "Error: Variable %s is constant. Cannot reassign a constant.\n", name)
        os.Exit(0)
        return nil
    }

    env.variables[name] = value

    if env.types[name] != value.VarType() && env.types[name] != ast.AnyType {
        if value.Type() == I8 || value.Type() == I16 || value.Type() == I32 || value.Type() == I64 || value.Type() == F32 || value.Type() == F64 {
            if env.types[name] == ast.Float64Type {
                return value
            }

            if i8Val, ok := value.(Int8Value); ok {
                if env.types[name] == ast.Int8Type && isInt8(float64(i8Val.Value)) {
                    return value
                }
            }

            if i16Val, ok := value.(Int16Value); ok {
                if env.types[name] == ast.Int16Type && isInt16(float64(i16Val.Value)) {
                    return value
                }
            }

            if i32Val, ok := value.(Int32Value); ok {
                if env.types[name] == ast.Int32Type && isInt32(float64(i32Val.Value)) {
                    return value
                }
            }

            if i64Val, ok := value.(Int64Value); ok {
                if env.types[name] == ast.Int64Type && isInt64(float64(i64Val.Value)) {
                    return value
                }
            }

            if f32Val, ok := value.(Float32Value); ok {
                if env.types[name] == ast.Float32Type && isFloat32(float64(f32Val.Value)) {
                    return value
                }
            }

            fmt.Fprintf(os.Stderr, "Error: Type mismatch, expected %s got %s\n", env.types[name], value.VarType())
            os.Exit(0)
        }
        fmt.Fprintf(os.Stderr, "Error: Type mismatch, expected %s got %s\n", env.types[name], value.VarType())
        os.Exit(0)
    }

    return value
}

func isInt8(value float64) bool {
    return value >= -0x80 && value <= 0x7F
}

func isInt16(value float64) bool {
    return value >= -0x8000 && value <= 0x7FFF
}

func isInt32(value float64) bool {
    return value >= -0x80000000 && value <= 0x7FFFFFFF
}

func isInt64(value float64) bool {
    return value >= -0x8000000000000000 && value <= 0x7FFFFFFFFFFFFFFF
}

func isFloat32(value float64) bool {
    return value >= -0x1.0p127 && value <= 0x1.0p127
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
        fmt.Fprintf(os.Stderr, "Error: %s not declared\n", name)
        os.Exit(0)
        return nil
    }

    return env.variables[name]
}

func (e *Environment) RemoveVariable(name string) {
    delete(e.variables, name)
}
