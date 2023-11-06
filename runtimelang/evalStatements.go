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

func EvaluateReturnStatement(statement ast.ReturnStatement, env Environment) RuntimeValue {
    returnValue, _ := Evaluate(statement.Value, env)
    return returnValue
}

func checkNumberTypes(value RuntimeValue, varType ast.VariableType) bool {
    return !(value.Type() == Number && varType == ast.Float64Type) && !(value.Type() == Number && varType == ast.Float32Type) && !(value.Type() == Number && varType == ast.Int64Type) && !(value.Type() == Number && varType == ast.Int32Type) && !(value.Type() == Number && varType == ast.Int16Type) && !(value.Type() == Number && varType == ast.Int8Type)
} 

func EvaluateVariableDeclaration(declaration ast.VariableDeclaration, env *Environment, varType ast.VariableType) RuntimeValue {
    value, _ := Evaluate(declaration.Value, *env)

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

    if value.Type() == Object && value.(ObjectValue).IsClass {
        value = value.(ObjectValue).Clone()
    }
    return env.DeclareVariable(declaration.Identifier, value, declaration.Constant, varType)
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
