package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlie/Jamlang/ast"
)

func EvaluateProgram(program ast.Program, env Environment) RuntimeValue {
    var lastEvaluated RuntimeValue = &InitialValue{}
    for _, statement := range program.Body {
        lastEvaluated, _ = Evaluate(statement, env)
    }
    return lastEvaluated
}

func EvaluateBreakStatement(statement ast.BreakStatement, env Environment) RuntimeValue {
    return &BreakType{}
}

func EvaluateContinueStatement(statement ast.ContinueStatement, env Environment) RuntimeValue {
    return &ContinueType{}
}

func EvaluateReturnStatement(statement ast.ReturnStatement, env Environment) RuntimeValue {
    returnValue, _ := Evaluate(statement.Value, env)
    return returnValue
}

func checkNumberTypes(value RuntimeValue, varType ast.VariableType) bool {
    return !(value.Type() == Number && varType == ast.Float64Type) && !(value.Type() == Number && varType == ast.Float32Type) && !(value.Type() == Number && varType == ast.Int64Type) && !(value.Type() == Number && varType == ast.Int32Type) && !(value.Type() == Number && varType == ast.Int16Type) && !(value.Type() == Number && varType == ast.Int8Type)
} 

func makeValueWithVarType(value RuntimeValue, varType ast.VariableType) RuntimeValue {
    switch varType {
    case ast.Float64Type:
        if _, ok := value.(FloatValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Float64Type, value.VarType())
            os.Exit(0)
        }
        return MakeFloat64Value(float64(value.(FloatValue).GetFloat()))
    case ast.Float32Type:
        if _, ok := value.(FloatValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Float32Type, value.VarType())
            os.Exit(0)
        }
        return MakeFloat32Value(float32(value.(FloatValue).GetFloat()))
    case ast.Int64Type:
        if _, ok := value.(IntValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Int64Type, value.VarType())
            os.Exit(0)
        }
        return MakeInt64Value(int64(value.(IntValue).GetInt()))
    case ast.Int32Type:
        if _, ok := value.(IntValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Int32Type, value.VarType())
            os.Exit(0)
        }
        return MakeInt32Value(int32(value.(IntValue).GetInt()))
    case ast.Int16Type:
        if _, ok := value.(IntValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Int16Type, value.VarType())
            os.Exit(0)
        }
        return MakeInt16Value(int16(value.(IntValue).GetInt()))
    case ast.Int8Type:
        if _, ok := value.(IntValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.Int8Type, value.VarType())
            os.Exit(0)
        }
        return MakeInt8Value(int8(value.(IntValue).GetInt()))
    case ast.ObjectType:
        if _, ok := value.(NullValue); ok {
            return value
        }
        if _, ok := value.(ObjectValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", ast.ObjectType, value.VarType())
            os.Exit(0)
        }
        return value
    case ast.NullType:
        return value
    case ast.AnyType:
        return value
    default:
        return value
    }
}

func EvaluateVariableDeclaration(declaration ast.VariableDeclaration, env *Environment, varType ast.VariableType) RuntimeValue {
    value, _ := Evaluate(declaration.Value, *env)

    actualValue := makeValueWithVarType(value, varType)

    if value.VarType() != varType && varType != ast.AnyType && !isNumber(value) {
        fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
        os.Exit(0)
    }

    if isNumber(value) {
        if declaration.Type == ast.Int8Type && value.VarType() != ast.Int8Type {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        } else if declaration.Type == ast.Int16Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type) {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        } else if declaration.Type == ast.Int32Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type && value.VarType() != ast.Int32Type) {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        } else if declaration.Type == ast.Int64Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type && value.VarType() != ast.Int32Type && value.VarType() != ast.Int64Type) {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        } else if declaration.Type == ast.Float32Type && value.VarType() != ast.Float32Type {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        } else if declaration.Type == ast.Float64Type && (value.VarType() != ast.Float32Type && value.VarType() != ast.Float64Type) {
            fmt.Fprintf(os.Stderr, "Error: Expected %s, got %s\n", declaration.Type, value.VarType())
            os.Exit(0)
        }
    }

    // if value.Type() == Object && value.(ObjectValue).IsClass {
    //     value = value.(ObjectValue).Clone()
    // }

    return env.DeclareVariable(declaration.Identifier, actualValue, declaration.Constant, varType)
}

func EvaluateVariableDeclarationDeprecated(declaration ast.VariableDeclaration, env *Environment) RuntimeValue {
    value, _ := Evaluate(declaration.Value, *env)
    if value.Type() == Object && value.(ObjectValue).IsClass {
        value = value.(ObjectValue).Clone()
    }
    return env.DeclareVariable(declaration.Identifier, value, declaration.Constant, ast.AnyType)
}

func EvaluateIdentifier(identifier *ast.Identifier, env *Environment) RuntimeValue {
    if identifier == nil {
        return MakeNullValue()
    }
    value := env.LookupVariable(identifier.Symbol)
    if value == nil {
        fmt.Fprintf(os.Stderr, "Undefined variable %s\n", identifier.Symbol)
        os.Exit(0)
        return nil
    }

    return value
}
